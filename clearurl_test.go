package clearurl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClear(t *testing.T) {
	h, err := Init()
	require.NoError(t, err)
	require.NoError(t, h.Preload())

	tests := map[string]struct {
		input string
		want  string
	}{
		"amazon": {input: "https://www.amazon.com/dp/exampleProduct/ref=sxin_0_pb?__mk_de_DE=ÅMÅŽÕÑ&keywords=tea&pd_rd_i=exampleProduct&pd_rd_r=8d39e4cd-1e4f-43db-b6e7-72e969a84aa5&pd_rd_w=1pcKM&pd_rd_wg=hYrNl&pf_rd_p=50bbfd25-5ef7-41a2-68d6-74d854b30e30&pf_rd_r=0GMWD0YYKA7XFGX55ADP&qid=1517757263&rnid=2914120011", want: "https://www.amazon.com/dp/exampleProduct"},
		"google": {input: "https://www.google.com/search?q=1&newwindow=1&iflsig=AHkkrS4AAAAAYk8jZ9ZArzEH0KuR&source=xx&spm=y7", want: "https://www.google.com/search?iflsig=AHkkrS4AAAAAYk8jZ9ZArzEH0KuR&newwindow=1&q=1"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := h.Clear(tc.input)
			require.NoError(t, err, name)
			require.Equal(t, tc.want, got, name)
		})
	}
}
