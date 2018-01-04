package commands

import (
	"errors"
	"strings"

	"github.com/Mentatious/mentat-cli/io"
	"github.com/Mentatious/mentat-cli/schema"
	"github.com/ybbus/jsonrpc"
)

// Delete ... delete selected entries
func Delete(rpcClient *jsonrpc.RPCClient, userid, UUIDS string, deleted *int) error {
	log := io.GetLog()
	uuidsList := strings.Split(UUIDS, ",")
	if (len(uuidsList) == 1) && (uuidsList[0] == "") {
		uuidsList = []string{}
	}
	rpcResp, err := rpcClient.CallNamed("entry.Delete", map[string]interface{}{
		"userid": userid,
		"uuids":  uuidsList,
	})
	if err != nil {
		log.Infof("[delete] error while deleting: %s", err.Error())
		return err
	}
	respData := schema.DeleteResponse{}
	err = rpcResp.GetObject(&respData)
	if err != nil {
		log.Infof("[delete] error unmarshalling response: %s", err)
		return err
	}
	if respData.Error != "" {
		log.Infof("[delete] receied error from API: %s", respData.Error)
		return errors.New(respData.Error)
	}
	*deleted = respData.Deleted
	return nil
}

// Cleanup ... cleanup DB
func Cleanup(rpcClient *jsonrpc.RPCClient, userid, types string, deleted *int) error {
	log := io.GetLog()
	typesList := strings.Split(types, ",")
	if (len(typesList) == 1) && (typesList[0] == "") {
		typesList = []string{}
	}
	rpcResp, err := rpcClient.CallNamed("entry.Cleanup", map[string]interface{}{
		"userid": userid,
		"types":  typesList,
	})
	if err != nil {
		log.Infof("[cleanup] error: %s", err.Error())
		return err
	}
	respData := schema.CleanupResponse{}
	err = rpcResp.GetObject(&respData)
	if err != nil {
		log.Infof("[cleanup] error unmarshalling response: %s", err)
		return err
	}
	if respData.Error != "" {
		log.Infof("[cleanup] receied error from API: %s", respData.Error)
		return errors.New(respData.Error)
	}
	*deleted = respData.Deleted
	return nil
}
