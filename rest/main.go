package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/irisnet/irishub/client/bank"
	"github.com/irisnet/irishub/app"

	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"time"
)

var seq int64

var seq2 int64

var cdc = app.MakeCodec()


func GetSequence(account string) int64{

	seq := int64(-1)

	res, body,err := Request( "1317", "GET", fmt.Sprintf("/bank/accounts/%s", account), nil)

	if res.StatusCode == http.StatusOK {

		var accInfo bank.BaseAccount
		//err = codec.Cdc.UnmarshalJSON([]byte(body), &resp)
		err = cdc.UnmarshalJSON([]byte(body), &accInfo)
		if err != nil{
			fmt.Println("error: ",err)
		} else{

			seq = accInfo.Sequence

		}
	}


	return seq
}


func GetAccount(account string)  bank.BaseAccount{

	var accInfo bank.BaseAccount

	res, body,err := Request( "1317", "GET", fmt.Sprintf("/bank/accounts/%s", account), nil)

	if res.StatusCode == http.StatusOK {

		//err = codec.Cdc.UnmarshalJSON([]byte(body), &resp)
		err = cdc.UnmarshalJSON([]byte(body), &accInfo)
		if err != nil {
			fmt.Println("error: ", err)
		}

	}

	return accInfo
}

func GetAccountByName(name string) keys.KeyOutput {
	var accInfo keys.KeyOutput


	res, body,err := Request( "1317", "GET", fmt.Sprintf("/keys/%s", name), nil)

	if res.StatusCode == http.StatusOK {

		//err = codec.Cdc.UnmarshalJSON([]byte(body), &resp)
		err = cdc.UnmarshalJSON([]byte(body), &accInfo)
		if err != nil {
			fmt.Println("error: ", err)
		}

	}

	return accInfo
}




func Request(port, method, path string, payload []byte) (*http.Response, string, error) {
	var (
		res *http.Response
	)
	url := fmt.Sprintf("http://104.198.127.16:%v%v", port, path)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))


	res, err = http.DefaultClient.Do(req)


	output, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return res, string(output),err
}

func SendTransaction(  port, name, password string, addr string) (receiveAddr sdk.AccAddress, resultTx ctypes.ResultBroadcastTxCommit){

	// send
	coinbz := sdk.NewInt64Coin("iris", 1).String()

	account := GetAccountByName(name)
	acc := GetAccount(account.Address.String())
	accnum := acc.AccountNumber
	sequence := acc.Sequence
	chainID := "fuxi-4000"

	fmt.Println(acc.AccountNumber, acc.Sequence,chainID)

	jsonStr := []byte(fmt.Sprintf(`{
		"base_tx":{
			"name":"%s",
			"password":"%s",
			"account_number":"%d",
			"sequence":"%d",
			"gas": "200000",
			"fee": "0.004iris",
			"chain_id":"%s"
        },
		"amount":"%s"
	}`, name, password, accnum, sequence, chainID, coinbz))
	res, body, _ := Request( port, "POST", fmt.Sprintf("/bank/%s/send", addr), jsonStr)

	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	fmt.Println(string(jsonStr))
	if res.StatusCode == http.StatusOK {

		fmt.Println("success")

	}else{
		fmt.Println("error",err)
	}

	return receiveAddr,resultTx
}
func SendTransactionWithSequence(  port, name, password string, addr string,seq int64) (receiveAddr sdk.AccAddress, resultTx ctypes.ResultBroadcastTxCommit){

	// send
	coinbz := sdk.NewInt64Coin("iris", 1).String()

	account := GetAccountByName(name)
	acc := GetAccount(account.Address.String())
	accnum := acc.AccountNumber
	chainID := "fuxi-4000"

	fmt.Println(acc.AccountNumber, acc.Sequence,chainID)

	jsonStr := []byte(fmt.Sprintf(`{
		"base_tx":{
			"name":"%s",
			"password":"%s",
			"account_number":"%d",
			"sequence":"%d",
			"gas": "200000",
			"fee": "0.004iris",
			"chain_id":"%s"
        },
		"amount":"%s"
	}`, name, password, accnum, seq, chainID, coinbz))
	res, body, _ := Request( port, "POST", fmt.Sprintf("/bank/%s/send?generate-only=false&async=true", addr), jsonStr)

	err := cdc.UnmarshalJSON([]byte(body), &resultTx)
	fmt.Println(string(jsonStr))
	if res.StatusCode == http.StatusOK {

		fmt.Println("success")

	}else{
		fmt.Println("error",err)
	}

	return receiveAddr,resultTx
}


func SendTwoTransactionWithSequence(  name_from string, name_to string, addr_from string,addr_to string,seq_from int64,seq_to int64) {

	// send

	_, result1 :=SendTransactionWithSequence("1317","iris","12345678",addr_to,seq_from)

	_, result2 :=SendTransactionWithSequence("1317","abc","12345678",addr_from,seq_to)

	fmt.Println(result1.Hash)
	fmt.Println(result2.Hash)

}



func main() {
	fmt.Println("Starting the application...")


	addr_from := "faa1u6mhz22ctc8t0j5fermakctl4n5tcsq53dh4xd"

	seq = GetSequence(addr_from)
	addr_to := "faa1n923ul8ckvudm30hghz7htg69aesjrgpzt7yfv"

	seq2= GetSequence(addr_to)

	fmt.Println("origin sequence is : ",seq)
	for t := range time.NewTicker(100 * time.Millisecond).C {
		heartBeat(t,seq)
		seq = seq +1
		seq2 = seq2 +1

	}

	fmt.Println("Terminating the application...")

}


func heartBeat(tick time.Time,seq int64){
	addr_from := "faa1u6mhz22ctc8t0j5fermakctl4n5tcsq53dh4xd"
	addr_to := "faa1n923ul8ckvudm30hghz7htg69aesjrgpzt7yfv"

	//
	//fmt.Println(seq)
	//_, result :=SendTransactionWithSequence("1317","iris","12345678",addr,seq)
	//

	SendTwoTransactionWithSequence("iris","abc",addr_from,addr_to,seq,seq2)
}
