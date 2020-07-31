import { BackendService, BackendMethod } from "./BackendService";
import { Subject } from "../Subject";


export interface ObjectV1
{
    key:            string;
    last_modified:  string;
    owner_id:       string;
    owner_name:     string;
    size:           number;
}


export class S3ObjectService
{
    private static _instance: S3ObjectService;


    public static getInstance ( ): S3ObjectService
    {
        if ( ! this._instance )
        {
            this._instance = new S3ObjectService();
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


    public async getObjects ( profileID: number, bucketName: string ): Promise<Array<ObjectV1>>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/profile/${encodeURIComponent(profileID)}/bucket/${encodeURIComponent(bucketName)}/object`,
            {
                method:     'GET',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        return resp.objects;
    }


    public async addObject ( profileID: number, bucketName: string, file: File, onProgress?: ( evt: any ) => any ): Promise<void>
    {
        const formData = new FormData();

        formData.set('file', file);

        await this._backendService.uploadFile(
            `/api/v1/profile/${encodeURIComponent(profileID)}/bucket/${encodeURIComponent(bucketName)}/object`,
            {
                method:     BackendMethod.POST,
                headers:    {
                    'Accept':       'application/json'
                },
                body: formData
            },
            onProgress
        );

        this._subjectUpdated.next();
    }


    public async deleteObject ( profileID: number, bucketName: string, objectKey: string ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/profile/${encodeURIComponent(profileID)}/bucket/${encodeURIComponent(bucketName)}/object/${encodeURIComponent(objectKey)}`,
            {
                method:     'DELETE',
                headers:    {
                    'Accept':       'application/json'
                },
            }
        );

        this._subjectUpdated.next();
    }



    public getDownloadObjectURI ( profileID: number, bucketName: string, objectKey: string ): string
    {
        return `/api/v1/profile/${encodeURIComponent(profileID)}/bucket/${encodeURIComponent(bucketName)}/object/${encodeURIComponent(objectKey)}`;
    }
}
