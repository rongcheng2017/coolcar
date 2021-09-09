import { rental } from "./proto_gen/rental/rental_pb";
import { Coolcar } from "./request";

export namespace ProfileService {
    export function getProfile(): Promise<rental.v1.IProfile> {
        return Coolcar.sendRequesWithAutyRetry(
            {
                method: 'GET',
                path: '/v1/profile',
                respMarshaller: rental.v1.Profile.fromObject,
            }
        )
    }

    export function submitProfile(req: rental.v1.IIdentity): Promise<rental.v1.IProfile> {
        return Coolcar.sendRequesWithAutyRetry({
            method: "POST",
            path: '/v1/profile',
            data: req,
            respMarshaller: rental.v1.Profile.fromObject
        })
    }

    export function clearProflie(): Promise<rental.v1.IProfile> {
        return Coolcar.sendRequesWithAutyRetry({
            method: 'DELETE',
            path: '/v1/profile',
            respMarshaller: rental.v1.Profile.fromObject,
        })
    }
}