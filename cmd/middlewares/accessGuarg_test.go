package middlewares

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Luckny/space-it/db/mock"
	db "github.com/Luckny/space-it/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRequireAccessLvl(t *testing.T) {
	user, _ := mockdb.RandomUser(t)
	space := mockdb.RandomSpace(t, user.ID)

	allPerms := mockdb.CreatePermission(t, user.ID, space.ID, true, true, true)
	nonePerms := mockdb.CreatePermission(t, user.ID, space.ID, false, false, false)

	testCases := []struct {
		name          string
		requestMethod string
		spaceID       string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:          "can read",
			requestMethod: http.MethodGet,
			spaceID:       space.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetPermissionsByUserAndSpaceIDParams{
					UserID:  user.ID,
					SpaceID: space.ID,
				}

				store.EXPECT().
					GetPermissionsByUserAndSpaceID(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(allPerms, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},

		{
			name:          "cannot read",
			requestMethod: http.MethodGet,
			spaceID:       space.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetPermissionsByUserAndSpaceIDParams{
					UserID:  user.ID,
					SpaceID: space.ID,
				}

				store.EXPECT().
					GetPermissionsByUserAndSpaceID(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nonePerms, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusForbidden)
			},
		},

		{
			name:          "can write",
			requestMethod: http.MethodPost,
			spaceID:       space.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetPermissionsByUserAndSpaceIDParams{
					UserID:  user.ID,
					SpaceID: space.ID,
				}

				store.EXPECT().
					GetPermissionsByUserAndSpaceID(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(allPerms, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},

		{
			name:          "cannot write",
			requestMethod: http.MethodPost,
			spaceID:       space.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetPermissionsByUserAndSpaceIDParams{
					UserID:  user.ID,
					SpaceID: space.ID,
				}

				store.EXPECT().
					GetPermissionsByUserAndSpaceID(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nonePerms, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusForbidden)
			},
		},

		{
			name:          "can delete",
			requestMethod: http.MethodDelete,
			spaceID:       space.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetPermissionsByUserAndSpaceIDParams{
					UserID:  user.ID,
					SpaceID: space.ID,
				}

				store.EXPECT().
					GetPermissionsByUserAndSpaceID(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(allPerms, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},

		{
			name:          "cannot delete",
			requestMethod: http.MethodDelete,
			spaceID:       space.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetPermissionsByUserAndSpaceIDParams{
					UserID:  user.ID,
					SpaceID: space.ID,
				}

				store.EXPECT().
					GetPermissionsByUserAndSpaceID(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nonePerms, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusForbidden)
			},
		},

		{
			name:          "admin permission",
			requestMethod: http.MethodPut,
			spaceID:       space.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetPermissionsByUserAndSpaceIDParams{
					UserID:  user.ID,
					SpaceID: space.ID,
				}

				store.EXPECT().
					GetPermissionsByUserAndSpaceID(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(allPerms, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},

		{
			name:          "invalid space id",
			requestMethod: http.MethodGet,
			spaceID:       "invalid-space-id",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetPermissionsByUserAndSpaceID(gomock.Any(), gomock.Any()).
					Times(0)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusInternalServerError)
			},
		},

		{
			name:          "sql error",
			requestMethod: http.MethodGet,
			spaceID:       space.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetPermissionsByUserAndSpaceIDParams{
					UserID:  user.ID,
					SpaceID: space.ID,
				}

				store.EXPECT().
					GetPermissionsByUserAndSpaceID(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Permission{}, sql.ErrConnDone)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusInternalServerError)
			},
		},

		{
			name:          "record not found",
			requestMethod: http.MethodGet,
			spaceID:       space.ID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetPermissionsByUserAndSpaceIDParams{
					UserID:  user.ID,
					SpaceID: space.ID,
				}

				store.EXPECT().
					GetPermissionsByUserAndSpaceID(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Permission{}, db.ErrRecordNotFound)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusForbidden)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init gomock
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// api server with mock store
			router := gin.Default()
			router.Use(func(c *gin.Context) {
				c.Set("user", &user)
				c.Next()
			})

			url := fmt.Sprintf("/spaces/%s/something", tc.spaceID)

			router.GET(
				"/spaces/:spaceID/something",
				RequireAccessLvl(ViewAccess, store),
				func(c *gin.Context) {
					c.JSON(http.StatusOK, nil)
				},
			)

			router.POST(
				"/spaces/:spaceID/something",
				RequireAccessLvl(WriteAccess, store),
				func(c *gin.Context) {
					c.JSON(http.StatusOK, nil)
				},
			)

			router.DELETE(
				"/spaces/:spaceID/something",
				RequireAccessLvl(DeleteAccess, store),
				func(c *gin.Context) {
					c.JSON(http.StatusOK, nil)
				},
			)

			router.PUT(
				"/spaces/:spaceID/something",
				RequireAccessLvl(AdminAccess, store),
				func(c *gin.Context) {
					c.JSON(http.StatusOK, nil)
				},
			)

			// create request
			request, err := http.NewRequest(tc.requestMethod, url, nil)
			require.NoError(t, err)

			// test recorder
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(recorder)
		})
	}
}
