// +build integration

package core

import (
	"context"
	"github.com/algorand/go-algorand-sdk/crypto"
	"testing"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/m2q/aema/core/client"
	"github.com/stretchr/testify/assert"
)

// Note: During integration tests, we need to make sure that there's no
// leak of goroutines. Because of this, every test will check for routine
// exits with WaitGroups returned by the SpawnManagingRoutine method of
// AlgorandBuffer

func TestIntegration_ValidAccount(t *testing.T) {
	buffer, err := CreateAlgorandBufferFromEnv()
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), client.AlgorandDefaultTimeout)
	info, err := buffer.Client.AccountInformation(buffer.AccountCrypt.Address.String(), ctx)
	cancel()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(info.CreatedApps))
}

// Test if app removal works
func TestIntegration_RemoveAccount(t *testing.T) {
	_ = createBufferAndRemoveApps(t)
}

// Remove application, and see if Manage re-creates the application
func TestIntegration_AccountGetsRestored(t *testing.T) {
	buffer := createBufferAndRemoveApps(t)
	wg, cancel := buffer.SpawnManagingRoutine()

	var info models.Account
	for !client.ValidAccount(info) {
		info, _ = buffer.Client.AccountInformation(buffer.AccountCrypt.Address.String(), context.Background())
		time.Sleep(time.Second)
	}

	cancel()
	if waitTimeout(wg, time.Second) {
		t.Fatalf("goroutine didn't finish in time")
	}
}

// Check if smart contract only allows deletion from an creator acc.
func TestSmartContract_DeleteWrongAcc(t *testing.T) {
	buffer, err := CreateAlgorandBufferFromEnv()
	assert.Nil(t, err)

	// Random Account deleting app should throw error
	randomAcc := crypto.GenerateAccount()
	err = buffer.Client.DeleteApplication(randomAcc, buffer.AppId)
	assert.NotNil(t, err)

	// Creator account deleting app should pass without problems
	err = buffer.Client.DeleteApplication(buffer.AccountCrypt, buffer.AppId)
	assert.Nil(t, err)
}

//func TestIntegration_PushData(t *testing.T) {
//	buffer, err := CreateAlgorandBufferFromEnv()
//	assert.Nil(t, err)
//
//	wg, cancel := buffer.SpawnManagingRoutine()
//
//	err = buffer.PutElements(map[string]string{
//		"554213" : "Astralis",
//	})
//	assert.Nil(t, err)
//
//	data, err := buffer.GetBuffer()
//	for len(data) == 0 {
//		assert.Nil(t, err)
//		time.Sleep(time.Millisecond * 200)
//		data, err = buffer.GetBuffer()
//	}
//	assert.EqualValues(t, 1, len(data))
//	t.Logf("data correctly inserted: [%s]", data)
//
//	// Make sure goroutine cancels in time
//	cancel()
//	if waitTimeout(wg, time.Second) {
//		t.Fatalf("goroutine didn't finish in time")
//	}
//}
