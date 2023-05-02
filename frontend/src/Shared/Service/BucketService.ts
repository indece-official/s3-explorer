import { BackendService } from './BackendService';
import { Subject } from 'ts-subject';


export interface BucketV1
{
    name: string;
}


export interface AddBucketV1Data
{
    name: string;
}


export interface BucketStatsV1
{
    files:      number;
    size:       number;
    complete:   boolean;
}


export class S3BucketService
{
    private static _instance: S3BucketService;


    public static getInstance ( ): S3BucketService
    {
        if ( ! this._instance )
        {
            this._instance = new S3BucketService();
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


    public async getBuckets ( profileID: number ): Promise<Array<BucketV1>>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/profile/${encodeURIComponent(profileID)}/bucket`,
            {
                method:     'GET',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        return resp.buckets;
    }
  
  
    public async getBucketStats ( profileID: number, bucketName: string, force: boolean ): Promise<BucketStatsV1>
    {
        const resp = await this._backendService.fetchJson(
            `/api/v1/profile/${encodeURIComponent(profileID)}/bucket/${encodeURIComponent(bucketName)}/stats?force=${encodeURIComponent(force)}`,
            {
                method:     'GET',
                headers:    {
                    'Accept':       'application/json'
                }
            }
        );

        return resp.stats;
    }


    public async addBucket ( profileID: number, data: AddBucketV1Data ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/profile/${encodeURIComponent(profileID)}/bucket`,
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
    }


    public async deleteBucket ( profileID: number, bucketName: string ): Promise<void>
    {
        await this._backendService.fetchJson(
            `/api/v1/profile/${encodeURIComponent(profileID)}/bucket/${encodeURIComponent(bucketName)}`,
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
