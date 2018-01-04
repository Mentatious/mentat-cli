package commands

import (
	"strings"

	"github.com/Mentatious/mentat-cli/io"
	"github.com/Mentatious/mentat-cli/schema"
	"github.com/ybbus/jsonrpc"
)

// Search ...
func Search(rpcClient *jsonrpc.RPCClient, userid, types, content, tags, priorities, query string) []schema.Entry {
	// TODO: implement raw queries support, until then "query" arg is ignored
	var results []schema.Entry
	log := io.GetLog()

	typesList := strings.Split(types, ",")
	if (len(typesList) == 1) && (typesList[0] == "") {
		typesList = []string{}
	}
	tagsList := strings.Split(tags, ",")
	if (len(tagsList) == 1) && (tagsList[0] == "") {
		tagsList = []string{}
	}
	prioList := strings.Split(priorities, ",")
	if (len(prioList) == 1) && (prioList[0] == "") {
		prioList = []string{}
	}

	rpcResp, err := rpcClient.CallNamed("entry.Search", map[string]interface{}{
		"userid":   userid,
		"types":    typesList,
		"content":  content,
		"tags":     tagsList,
		"priority": prioList,
	})
	if err != nil {
		log.Infof("[search] error while searching: %s", err.Error())
		return results
	}
	respData := schema.SearchResponse{}
	err = rpcResp.GetObject(&respData)
	if err != nil {
		log.Infof("[search] error unmarshalling response: %s", err)
		return results
	}
	return respData.Entries
}
