import * as React from 'react';
import { DownloadManagerService, DownloadStatus } from './DownloadManagerService';
import { ProgressBar } from '../ProgressBar/ProgressBar';

import './DownloadManager.css';


export interface DownloadManagerProps
{
    onError: ( error: Error | null ) => any;
}


interface DownloadManagerState
{
    downloads: Array<DownloadStatus>;
}


export class DownloadManager extends React.Component<DownloadManagerProps, DownloadManagerState>
{
    private readonly _downloadManagerService: DownloadManagerService;


    constructor ( props: DownloadManagerProps )
    {
        super(props);

        this.state = {
            downloads: []
        };

        this._downloadManagerService = DownloadManagerService.getInstance();

        this._load = this._load.bind(this);
    }


    private _load ( download?: DownloadStatus ): void
    {
        this.setState({
            downloads: this._downloadManagerService.getActiveDownloads()
        });

        if ( download && download.finished && download.error )
        {
            this.props.onError(download.error);
        }
    }


    public componentDidMount ( ): void
    {
        this._load();

        this._downloadManagerService.updated().subscribe(this, this._load);
    }


    public componentWillUnmount ( ): void
    {
        this._downloadManagerService.updated().unsubscribe(this);
    }


    public render ( )
    {
        if ( this.state.downloads.length === 0 )
        {
            return null;
        }

        return (
            <div className='DownloadManager'>
                <div className='DownloadManager-downloads'>
                    {this.state.downloads.map( ( dl ) =>
                        <div
                            className='DownloadManager-download'
                            key={dl.uid}>
                            <div className='DownloadManager-download-filename'>
                                {dl.object_key}
                            </div>

                            <ProgressBar
                                value={dl.loaded}
                                total={dl.total}
                            />
                        </div>
                    )}
                </div>
            </div>
        );
    }
}
