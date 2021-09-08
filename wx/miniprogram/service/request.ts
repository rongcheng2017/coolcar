import camelcaseKeys = require("camelcase-keys")
import { auth } from "./proto_gen/auth/auth_pb"

const bearerPrefix = "Bearer "

export namespace Coolcar {
    const serverAddr = 'http://localhost:8080'
    const AUTH_ERR = 'AUTH_ERR'

    const authData = {
        token: '',
        expiryMs: 0,
    }

    interface RequestOption<REQ, RES> {
        method: 'GET' | 'PUT' | 'POST' | 'DELETE'
        path: string
        data: REQ
        respMarshaller: (r: object) => RES

    }
    export interface AuthOption {
        attachAuthHeader: boolean,
        retryOnAuthErr: boolean,
    }

    export async function sendRequesWithAutyRetry<REQ, RES>(o: RequestOption<REQ, RES>, a?: AuthOption): Promise<RES> {
        const authOpt = a || {
            attachAuthHeader: true,
            retryOnAuthErr: true,
        }
        try {
            await login()
            return sendRequest(o, a!)
        } catch (error) {
            if (error === AUTH_ERR && authOpt.retryOnAuthErr) {
                authData.token = ''
                authData.expiryMs = 0
                return sendRequesWithAutyRetry(o, {
                    attachAuthHeader: authOpt.attachAuthHeader,
                    retryOnAuthErr: false
                })
            } else {
                throw error
            }
        }

    }
    export async function login() {

        if (authData.token && authData.expiryMs >= Date.now()) {
            //如果token有效
            return
        }

        const wxResp = await wxLogin()
        const reqTimeMs = Date.now()
        const resp = await sendRequest<auth.v1.ILoginRequest, auth.v1.ILoginResponse>({
            method: 'POST',
            path: '/v1/auth/login',
            data: {
                code: wxResp.code
            },
            respMarshaller: auth.v1.LoginResponse.fromObject
        }, {
            attachAuthHeader: false,
            retryOnAuthErr: false,
        })
        authData.token = resp.accessToken!
        authData.expiryMs = reqTimeMs + resp.expiresIn! * 1000
    }

    function sendRequest<REQ, RES>(o: RequestOption<REQ, RES>, a: AuthOption): Promise<RES> {
        const authOpt = a || {
            attachAuthHeader: true
        }
        return new Promise((resolve, reject) => {
            const header: Record<string, any> = {}
            if (authOpt.attachAuthHeader) {
                if (authData.token && authData.expiryMs >= Date.now()) {
                    header.authorization = bearerPrefix + authData.token
                } else {
                    reject(AUTH_ERR)
                    return
                }
            }
            wx.request(
                {
                    url: serverAddr + o.path,
                    method: o.method,
                    data: o.data,
                    header,
                    success: res => {
                        if (res.statusCode === 401) {
                            reject(AUTH_ERR)
                        } else if (res.statusCode >= 400) {
                            reject(res)
                        } else {
                            resolve(o.respMarshaller(camelcaseKeys(res.data as object, { deep: true })))
                        }
                    },
                    fail: reject
                }
            )
        })
    }

    function wxLogin(): Promise<WechatMiniprogram.LoginSuccessCallbackResult> {

        return new Promise((resolve, reject) => {
            wx.login({
                success: resolve,
                fail: reject
            })
        })
    }
}