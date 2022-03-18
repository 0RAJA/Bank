package token

//
//func TestJwtMaker(t *testing.T) {
//	maker, err := NewJWTMaker(utils.RandomString(32))
//	require.NoError(t, err)
//
//	username := utils.RandomString(10)
//	duration := time.Minute
//	issuedAt := time.Now()
//	expiredAt := issuedAt.Add(duration)
//	token, err := maker.CreateToken(username, duration)
//	require.NoError(t, err)
//	require.NotEmpty(t, token)
//
//	payload, err := maker.VerifyToken(token)
//	require.NoError(t, err)
//	require.NotEmpty(t, payload)
//
//	require.Equal(t, payload.UserName, username)
//	require.Equal(t, payload.IssuedAt, issuedAt)
//
//	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
//}
