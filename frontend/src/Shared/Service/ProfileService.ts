import { BackendService } from './BackendService';
import { Subject } from 'ts-subject';


export interface ProfileV1
{
    id:         number;
    name:       string;
    access_key: string;
    secret_key: string;
    region:     string;
    endpoint:   string;
    ssl:        boolean;
    path_style: boolean;
    buckets:    Array<string>;
}


export interface AddProfileV1Data
{
    name:       string;
    access_key: string;
    secret_key: string;
    region:     string;
    endpoint:   string;
    ssl:        boolean;
    path_style: boolean;
    buckets:    Array<string>;
}


export interface UpdateProfileV1Data
{
    name:       string;
    access_key: string | null;
    secret_key: string | null;
    region:     string;
    endpoint:   string;
    ssl:        boolean;
    path_style: boolean;
    buckets:    Array<string>;
}


export class S3ProfileService
{
    private static _instance: S3ProfileService;


    public static getInstance ( ): S3ProfileService
    {
        if ( ! this._instance )
        {
            this._instance = new S3ProfileService();
        }

        return this._instance;
    }


    private readonly _backendService: BackendService;
    private readonly _subjectUpdated: Subject<void>;


    constructor ( )
    {
        this._backendService = BackendService.getInstance();
        this._subjectUpdated = new Subject();
    }


    public updated ( ): Subject<void>
    {
        return this._subjectUpdated;
    }


    public async getProfiles ( ): Promise<Array<ProfileV1>>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/profile`,
            {
                method:     'GET',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        return resp.profiles;
    }
    
    
    public async addProfile ( data: AddProfileV1Data ): Promise<number>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/profile`,
            {
                method:     'POST',
                headers:    {
                    'Accept':       'application/json',
                    'Content-type': 'application/json'
                },
                body: JSON.stringify(data)
            }
        );

        this._subjectUpdated.next();

        return resp.profile_id;
    }


    public async updateProfile ( profileID: number, data: UpdateProfileV1Data ): Promise<number>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/profile/${encodeURIComponent(profileID)}`,
            {
                method:     'PUT',
                headers:    {
                    'Accept':       'application/json',
                    'Content-type': 'application/json'
                },
                body: JSON.stringify(data)
            }
        );

        this._subjectUpdated.next();

        return resp.profile_id;
    }


    public async deleteProfile ( profileID: number ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/profile/${encodeURIComponent(profileID)}`,
            {
                method:     'DELETE',
                headers:    {
                    'Accept':       'application/json'
                },
            }
        );

        this._subjectUpdated.next();
    }
}
