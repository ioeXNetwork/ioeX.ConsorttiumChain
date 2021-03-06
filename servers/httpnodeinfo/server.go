package httpnodeinfo

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	chain "github.com/ioeX/ioeX.SideChain/blockchain"
	"github.com/ioeX/ioeX.SideChain/config"
	"github.com/ioeX/ioeX.SideChain/servers"
)

type Info struct {
	NodeVersion   string
	BlockHeight   uint32
	NeighborCnt   int
	Neighbors     []NgbNodeInfo
	HttpRestPort  int
	HttpWsPort    int
	HttpJsonPort  int
	HttpLocalPort int
	NodePort      uint16
	NodeId        string
}

type NgbNodeInfo struct {
	NgbId   string
	NbrAddr string
}

var node = servers.NodeForServers

var templates = template.Must(template.New("info").Parse(page))

func viewHandler(w http.ResponseWriter, r *http.Request) {
	var ngbrNodersInfo []NgbNodeInfo
	ngbrNoders := node.GetNeighborNoder()

	for i := 0; i < len(ngbrNoders); i++ {
		ngbrNodersInfo = append(ngbrNodersInfo, NgbNodeInfo{
			NgbId:   fmt.Sprintf("0x%x", ngbrNoders[i].ID()),
			NbrAddr: ngbrNoders[i].Addr() + ":" + strconv.Itoa(ngbrNoders[i].HttpInfoPort()),
		})
	}

	pageInfo := &Info{
		BlockHeight:  chain.DefaultLedger.Blockchain.BlockHeight,
		NeighborCnt:  len(ngbrNoders),
		Neighbors:    ngbrNodersInfo,
		HttpRestPort: config.Parameters.HttpRestPort,
		HttpWsPort:   config.Parameters.HttpWsPort,
		HttpJsonPort: config.Parameters.HttpJsonPort,
		NodePort:     config.Parameters.NodePort,
		NodeId:       fmt.Sprintf("0x%x", node.ID()),
	}

	err := templates.ExecuteTemplate(w, "info", pageInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func StartServer() {
	http.HandleFunc("/info", viewHandler)
	http.ListenAndServe(":"+strconv.Itoa(int(config.Parameters.HttpInfoPort)), nil)
}
