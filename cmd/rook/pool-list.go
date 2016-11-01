package rook

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/rook/rook/pkg/model"
	"github.com/rook/rook/pkg/rook/client"
	"github.com/rook/rook/pkg/util/display"
	"github.com/rook/rook/pkg/util/flags"
	"github.com/spf13/cobra"
)

var poolListCmd = &cobra.Command{
	Use:   "ls",
	Short: "Gets a listing with details of all storage pools in the cluster",
}

func init() {
	poolListCmd.RunE = listPoolsEntry
}

func listPoolsEntry(cmd *cobra.Command, args []string) error {
	if err := flags.VerifyRequiredFlags(cmd, []string{}); err != nil {
		return err
	}

	c := client.NewRookNetworkRestClient(client.GetRestURL(apiServerEndpoint), http.DefaultClient)
	out, err := listPools(c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Print(out)
	return nil
}

func listPools(c client.RookRestClient) (string, error) {
	pools, err := c.GetPools()
	if err != nil {
		return "", fmt.Errorf("failed to get pools: %+v", err)
	}

	var buffer bytes.Buffer
	w := NewTableWriter(&buffer)

	fmt.Fprintln(w, "NAME\tNUMBER\tTYPE\tSIZE\tDATA\tCODING\tALGORITHM")

	for _, p := range pools {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s\t%s\n", p.Name, p.Number, model.PoolTypeToString(p.Type),
			display.NumToStrOmitEmpty(p.ReplicationConfig.Size),
			display.NumToStrOmitEmpty(p.ErasureCodedConfig.DataChunkCount),
			display.NumToStrOmitEmpty(p.ErasureCodedConfig.CodingChunkCount),
			p.ErasureCodedConfig.Algorithm)
	}

	w.Flush()
	return buffer.String(), nil
}