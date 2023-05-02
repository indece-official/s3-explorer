import { Subject } from 'ts-subject';
import { S3ObjectService } from "../Service/ObjectService";

export interface DownloadStatus
{
    uid:            string;
    profile_id:     number;
    bucket_name:    string;
    object_key:     string;
    loaded:         number;
    total:          number;
    finished:       boolean;
    error:          Error | null;
}


export class DownloadManagerService
{
    private static _instance: DownloadManagerService;


    public static getInstance ( ): DownloadManagerService
    {
        if ( ! this._instance )
        {
            this._instance = new DownloadManagerService();
        }

        return this._instance;
    }


    private readonly _s3ObjectService: S3ObjectService;
    private readonly _subjectUpdated: Subject<DownloadStatus>;
    private readonly _downloads: {[key: string]: DownloadStatus};


    constructor ( )
    {
        this._s3ObjectService = S3ObjectService.getInstance();
        this._subjectUpdated = new Subject();
        this._downloads = {};
    }


    public updated ( ): Subject<DownloadStatus>
    {
        return this._subjectUpdated;
    }


    public getDownloads ( ): Array<DownloadStatus>
    {
        return Object.values(this._downloads);
    }


    public getActiveDownloads ( ): Array<DownloadStatus>
    {
        return Object.values(this._downloads)
            .filter( ( dl ) => !dl.finished );
    }



    public async download ( profileID: number,
                            bucketName: string,
                            objectKey: string ): Promise<void>
    {
        if ( typeof((window as any).s3DownloadFile) === 'function' )
        {
            const uid = await (window as any).s3DownloadFile(
                profileID,
                bucketName,
                objectKey
            );

            this._downloads[uid] = await (window as any).s3DownloadStatus(uid);

            this._subjectUpdated.next(this._downloads[uid]);

            const interval = setInterval( async ( ) =>
            {
                this._downloads[uid] = await (window as any).s3DownloadStatus(uid)

                this._subjectUpdated.next(this._downloads[uid]);

                if ( this._downloads[uid].finished )
                {
                    clearInterval(interval);
                }
            }, 250);
        }
        else
        {
            const url = this._s3ObjectService.getDownloadObjectURI(
                profileID,
                bucketName,
                objectKey
            );

            window.open(url, '_blank');
        }
    }
}
