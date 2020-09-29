package authprovider

import (
	"fmt"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHydra_ClientWithIDAndSecret(t *testing.T) {
	h := Hydra{}
	c := h.ClientWithIDAndSecret("", "")

	resp, err := c.Get("http://localhost:4444/userinfo")
	require.NoError(t, err)

	b, _ := httputil.DumpResponse(resp, true)
	fmt.Println(string(b))
	t.FailNow()
}
