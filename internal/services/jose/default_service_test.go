package jose

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeusito/gochat/config"
)

var logger = config.NewLogger()

func TestNoPublicKeyPassed(t *testing.T) {
	service, err := NewService(logger, "")

	assert.Empty(t, service, "Service should've been empty")
	assert.Errorf(t, err, "Invalid error thrown")
}

func TestKeysNotCorrectlyEncoded(t *testing.T) {
	service, err := NewService(logger, "dsdasdasdzczxcxcxc")

	assert.Empty(t, service, "Service should've been empty")
	assert.Errorf(t, err, "Invalid error thrown")
}

func TestValidateExpiredToken(t *testing.T) {
	pubKey := "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUFrYk5DVGR5OHREZzBHLzUvM2ZldAovK1BJeDdZeUVSVktOWGxtTHZQdTFQTGxMd2ppNDJQSGF5aVBqUU5jUVVtZXp2Y25NYjl2Z0Y4ZEF0OVROUU9qCkNxYlNLdnV4eGM0dWJxMSs1NFpkbEJ4VjliWmJkWmlpODU2RFJqdXVSN1c3cUpXbTVGZEVqdUNxWCsySjIxNWwKbE1HU29hblgzWTlSNDRRS1hFYjhjaG83OGFTZkFid2JBeC9DSVhNZkZGUkk5aXNPQWdTbitMNkkyM0IrMGIxbgpPeTY1Wjh6RVp0eVBJR1ZIbEtzWHdUOEhERG1LTHlvSENiWEhOcUdmQ1pYWFlPeFB1SXNvbmcwZCt4d24wZFBECnRvTDJwTFBkL2tvb200Tmo4dlQvVTRxRzNSMk5qWk9GRmx1elg5RDNWWUhpZ01LbytXUCthMkFKTjJNQ0dtQUsKWmRadC96WFNTK2JyK2RneEFxMm5Edk1sSXVxYzBBUEFFT21LTVhvS3RUVm9tWng0MnF4RTlyQlI4Q3lPek9jSApLU3FNUWhZSEQ4dDVYNVdRU0NqY0pmUzF1QmdNYTNrS3VoNDlncEgycGpTODQ4eFljNUFVR0VUUUliMENnaWlHCkZkZDdRQVNBaGV0dmRCMXhLaTBqWC9VMGVPdmUzM0loOG5FT21qbWZGNWR0MUZZKzF3RGNzNzArSStXR0dvRC8KQ0k5MlRucTB3VWhXNFE3cGhtVlZUN1QxVGV4UWZ3S011MWxqYmYzdzJxSHIxcEpGQkRBWUVJTmZXQi9qNTFNSgpYbXZsNjVzRXIwUk00dXprY0hjcmdOei9yeTlVQTFqNEhBTGF4YmFuUk1sQVUvc2Q1MkN3M1cwendaTDZ2bmU3CklqTmpWSEp1T2FFcmtRckpkeVYvdng4Q0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ=="
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyTmFtZSI6InRlc3QiLCJ1c2VySWQiOiJ0ZXN0IiwiaXNzIjoiemV1c2l0byIsInN1YiI6InRlc3QiLCJhdWQiOlsiemV1c2l0by5tZSJdLCJleHAiOjE2MzgzNzI5MTEsIm5iZiI6MTYzODM2OTMxMSwiaWF0IjoxNjM4MzY5MzExfQ.Ks2RiRYNpd7zGa2cnQpA3HW8RIAs3a1sUcB3pcldxWJjLX0viDDQaoVvzGnvrg6tTDKok2e89i-tUDCzTYtYE3tZyTI_sRpnBJpCcy2aUMLrDojU3aVMF-GajNq2X1j_OQuiWlFEI4REeiyMEdEpR67MRdchUuYWWPI3X-Aep0VKAyMo_lAx2q5l2wVhCGQr7OXrz3ZwiVAXQtaF_RTc1KKbzqvxzcxxNeidWxuAKjaaXsrGOrVGcQi6W8qrBp-C26ERUjOk2TClz_IWN6Wr1tRblTM-jWTlIs89lyZqLeZq-C2yNRsMgVyPL9_q6wc0wwyfb0Uml3Bd9m_6dBdtD9QFcbtiUUPtw4JWrVD15JBDnhkApM6M1AnjMp719vFm878QWMzQOwMHEp8EI8KCwOJACNSGqpUWU2zP7GxFVwkvbgc_qRAlxNJxMTsHprL1_vCmT31-GyIwX4RvolVZIGoCQ79uIooJfgJAQqlMEyo31ykypVCptjdKzeMSD1qYbors-GWt2suUcsbjGuPThhlJa6XM0HNjV4sXuWrSA9gBhDcaAN-sXb-WF3QK6Yo2Okg6uxAYUC6vKAjx_vjOSn-Bc6EFCcSGKFSRcMZD3NzD6ciRA70iOIA209YWtSgmCsnyW4yyPWEHMWYsYVqNFuVf9Zn9GgJa0ICAJw75wSQ"

	service, err := NewService(logger, pubKey)

	assert.NoError(t, err, "No error should've been thrown")

	_, err = service.Parse(token)

	assert.Errorf(t, err, "Token is expired")
}

func TestValidateTokenSignature(t *testing.T) {
	pubKey := "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUF0VkNQSHllOVprYXU3K2lVY2dJMwpWZlF1Q21weVE1L2NaSk9YY281YWJYR0Q2M3cvTWk5RHhDK3NnVUIvUGVjMWxMRGhGZVhkQ2dNaWNVeHpKZFNpClFDcTRZL0E4aFgxN0FackUwWXR5VDFjd0FvL09ZRXYvU0pySWFtcFErOUVnS3JSUGUvNjlUcVREalhvZ1c2THkKY0tRNTBGYk1YN2ZZMWZCazJGd1BRa3VQR1h6NnFDRUo2VnU1empPL01FdStpWkYzNStZWE1WanpRa2NmY3dqbgpla1QrUmtTMVFnNHNpdWhEaE1zekpvUk85cVo2VmFtRWFBbDk2c0pKUytIQUtReElRM1ROcnpvdEFGVm55UTVICkl1RHJMWFh6ZWlGN0VpeVhkR0Q0TW5malcwTis5UnlDcmNLa2U5Qit5VERPZjJ5aVFaUEJtU20zcUhjQ05PdXgKTTUyWlZLRzVETGlkWnhQdFdEc2l5RW9rdnZWMDliVjdPL0xYOTVMWEg2V1NGS0llS3NuVUVMdGZKbWNiMVJsZwo0dkV0RDVLQndaNTBxL1RxUXZYYXZlYkxVUXNEUnk1akFJWkRRWFl2Z0tEK0tOMjhpNmsrMWRsWHkxL2JKQnQzCkpWWkhiTXVVV2FvOEFWL08wMjF0WmxCNDJCTWs3Mm90SVg3STMxa3U3QjJpS21tbkY3MDd3bFdUR3VjbmFyREwKTVk5dnFDNkR3cklPVEdXK09ERVROU0ZVTUdvekxML3NTSDk2TlFCUCtvT3Y4WUVjUVdNMU5VY21QajhPTmdSbQp6V3MxV0IyaG9BUUJmdWF2K3pySzd0VTA0RXVzOFNaYkY4K3lubmJERmFiT3YxbHJDYUs3cnU1a3JNQWwwU2tmCkF2YWR5dHRGemFzOFBkZkhxY0huZSs4Q0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ=="
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyTmFtZSI6InRlc3QiLCJ1c2VySWQiOiJ0ZXN0IiwiaXNzIjoiemV1c2l0byIsInN1YiI6InRlc3QiLCJhdWQiOlsiemV1c2l0by5tZSJdLCJleHAiOjE3MzgzNzI5MTEsIm5iZiI6MTYzODM2OTMxMSwiaWF0IjoxNjM4MzY5MzExfQ.jZ8X5A4x_S0itXvCQ3MoCJPojrW_wSOsarA1KG1_iCp3_eGSM62qcKpSgns6XNrTb0CDiVqxm2QnWunFEVb04MnqZyA5i42YXW_vOMj-jdcFg4RfVwKvN9RKV5OA1hYXhb13dTux7azPHstD4YUw7xgt_D6N3aDk6-zO8fYOTaRPuMp2aQu1P4191dO28hZF0W_RmluW4RGvx48a-l4KApC4IJ66UBYve0TRXdxlzDcFOK0QPs36jDDgVxvsG-fZ-F-imcGLFyYdXzX9gDc0Em2OSOTxWvBVe_UNBVyFjTmEpWwfQfxNGTmE0sMuU6GWNEhA7-EPKMxQeBT9c8REzSZPNBwLRYVsK5mUKRtU-x3BRvKpEF2QY6rE4DCChEUJ7bF-Q_dqUjJj_G3rfO9hxFMeDE4MgbFvbfEOIHv7Gp37KNEew0v22wfr9vkNN2VoVzEShanm8wlTwgSqByOq_wU8GUIAwAjotSmTy9puLG8SlzHZoWg0ylLwieCnBMh9MIt6Kh8vZMNiWVmGucMZ9J2ntGf8iekvTrOIsAfpJmA9x7u3quvRQYs8loLB7ZB2MexjN2_NXeRjeO3kbSKwPcv7xlRXN7pjWYgP2kZCROzn3Mmk7K7hH_W_o1dCpGu-84dsb7veztpV_7DbHIKklFhnZjfOhVsCyKVGM8fu25k"
	service, err := NewService(logger, pubKey)

	assert.NoError(t, err, "No error should've been thrown")

	_, err = service.Parse(token)

	assert.Errorf(t, err, "crypto/rsa: verification error")
}

func TestCorrectParsing(t *testing.T) {
	pubKey := "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUFrYk5DVGR5OHREZzBHLzUvM2ZldAovK1BJeDdZeUVSVktOWGxtTHZQdTFQTGxMd2ppNDJQSGF5aVBqUU5jUVVtZXp2Y25NYjl2Z0Y4ZEF0OVROUU9qCkNxYlNLdnV4eGM0dWJxMSs1NFpkbEJ4VjliWmJkWmlpODU2RFJqdXVSN1c3cUpXbTVGZEVqdUNxWCsySjIxNWwKbE1HU29hblgzWTlSNDRRS1hFYjhjaG83OGFTZkFid2JBeC9DSVhNZkZGUkk5aXNPQWdTbitMNkkyM0IrMGIxbgpPeTY1Wjh6RVp0eVBJR1ZIbEtzWHdUOEhERG1LTHlvSENiWEhOcUdmQ1pYWFlPeFB1SXNvbmcwZCt4d24wZFBECnRvTDJwTFBkL2tvb200Tmo4dlQvVTRxRzNSMk5qWk9GRmx1elg5RDNWWUhpZ01LbytXUCthMkFKTjJNQ0dtQUsKWmRadC96WFNTK2JyK2RneEFxMm5Edk1sSXVxYzBBUEFFT21LTVhvS3RUVm9tWng0MnF4RTlyQlI4Q3lPek9jSApLU3FNUWhZSEQ4dDVYNVdRU0NqY0pmUzF1QmdNYTNrS3VoNDlncEgycGpTODQ4eFljNUFVR0VUUUliMENnaWlHCkZkZDdRQVNBaGV0dmRCMXhLaTBqWC9VMGVPdmUzM0loOG5FT21qbWZGNWR0MUZZKzF3RGNzNzArSStXR0dvRC8KQ0k5MlRucTB3VWhXNFE3cGhtVlZUN1QxVGV4UWZ3S011MWxqYmYzdzJxSHIxcEpGQkRBWUVJTmZXQi9qNTFNSgpYbXZsNjVzRXIwUk00dXprY0hjcmdOei9yeTlVQTFqNEhBTGF4YmFuUk1sQVUvc2Q1MkN3M1cwendaTDZ2bmU3CklqTmpWSEp1T2FFcmtRckpkeVYvdng4Q0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ=="
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyTmFtZSI6InRlc3QiLCJ1c2VySWQiOiJ0ZXN0IiwiaXNzIjoiemV1c2l0byIsInN1YiI6InRlc3QiLCJhdWQiOlsiemV1c2l0by5tZSJdLCJleHAiOjE3MzgzNzI5MTEsIm5iZiI6MTYzODM2OTMxMSwiaWF0IjoxNjM4MzY5MzExfQ.jZ8X5A4x_S0itXvCQ3MoCJPojrW_wSOsarA1KG1_iCp3_eGSM62qcKpSgns6XNrTb0CDiVqxm2QnWunFEVb04MnqZyA5i42YXW_vOMj-jdcFg4RfVwKvN9RKV5OA1hYXhb13dTux7azPHstD4YUw7xgt_D6N3aDk6-zO8fYOTaRPuMp2aQu1P4191dO28hZF0W_RmluW4RGvx48a-l4KApC4IJ66UBYve0TRXdxlzDcFOK0QPs36jDDgVxvsG-fZ-F-imcGLFyYdXzX9gDc0Em2OSOTxWvBVe_UNBVyFjTmEpWwfQfxNGTmE0sMuU6GWNEhA7-EPKMxQeBT9c8REzSZPNBwLRYVsK5mUKRtU-x3BRvKpEF2QY6rE4DCChEUJ7bF-Q_dqUjJj_G3rfO9hxFMeDE4MgbFvbfEOIHv7Gp37KNEew0v22wfr9vkNN2VoVzEShanm8wlTwgSqByOq_wU8GUIAwAjotSmTy9puLG8SlzHZoWg0ylLwieCnBMh9MIt6Kh8vZMNiWVmGucMZ9J2ntGf8iekvTrOIsAfpJmA9x7u3quvRQYs8loLB7ZB2MexjN2_NXeRjeO3kbSKwPcv7xlRXN7pjWYgP2kZCROzn3Mmk7K7hH_W_o1dCpGu-84dsb7veztpV_7DbHIKklFhnZjfOhVsCyKVGM8fu25k"
	service, err := NewService(logger, pubKey)

	assert.NoError(t, err, "No error should've been thrown")

	claims, err := service.Parse(token)

	assert.NoError(t, err, "No error should've been thrown")
	assert.Equal(t, "test", claims.UserID, "id should be test")
	assert.Equal(t, "test", claims.UserName, "name should be test")
}