import * as React from 'react';
import Bytes from 'bytes';
import Moment from 'moment';
import InfiniteScroll from 'react-infinite-scroller';
import { ObjectV1, S3ObjectService } from '../Service/ObjectService';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSync, faTimes, faDownload } from '@fortawesome/free-solid-svg-icons';
import { DownloadManagerService } from '../DownloadManager/DownloadManagerService';
import { Spinner } from '../Spinner/Spinner';
import { S3ProfileService } from '../Service/ProfileService';

import './ObjectList.css';


export interface ObjectListProps
{
    profileID:      number;
    bucketName:     string;
    onSelectObject: ( object: ObjectV1 ) => any;
    onDeleteObject: ( object: ObjectV1 ) => any;
    onError:        ( err: Error | null ) => any;
}


interface ObjectListState
{
    objects:            Array<ObjectV1>;
    hasMore:            boolean;
    loading:            boolean;
    continuationToken:  string;
}


export class ObjectList extends React.Component<ObjectListProps, ObjectListState>
{
    private readonly BULK_SIZE                  = 100;
    private readonly _s3ProfileService:         S3ProfileService;
    private readonly _s3ObjectService:          S3ObjectService;
    private readonly _downloadManagerService:   DownloadManagerService;


    constructor ( props: ObjectListProps )
    {
        super(props);

        this.state = {
            objects:            [],
            hasMore:            false,
            loading:            false,
            continuationToken:  ''
        };

        this._s3ProfileService = S3ProfileService.getInstance();
        this._s3ObjectService = S3ObjectService.getInstance();
        this._downloadManagerService = DownloadManagerService.getInstance();

        this._load = this._load.bind(this);
        this._loadMore = this._loadMore.bind(this);
    }


    private async _load ( ): Promise<void>
    {
        if ( !this.props.profileID || !this.props.bucketName )
        {
            this.setState({
                objects:            [],
                continuationToken:  '',
                hasMore:            false
            });

            return;
        }

        this.setState({
            loading:    true,
            objects:    []
        });

        try
        {
            const resp = await this._s3ObjectService.getObjects(
                this.props.profileID,
                this.props.bucketName,
                "",
                this.BULK_SIZE
            );

            this.setState({
                objects:            resp.objects,
                continuationToken:  resp.continuation_token,
                hasMore:            !!(resp.continuation_token && resp.objects.length >= this.BULK_SIZE),
                loading:            false
            });

            this.props.onError(null);
        }
        catch ( err )
        {
            console.error(`Error loading objects: ${(err as Error).message}`, err);

            this.props.onError(err as Error);

            this.setState({
                loading:    false
            });
        }
    }
    
    private async _loadMore ( ): Promise<void>
    {
        if ( !this.state.continuationToken || !this.state.hasMore || this.state.loading )
        {
            return;
        }

        this.setState({
            loading:    true
        });

        try
        {
            const resp = await this._s3ObjectService.getObjects(
                this.props.profileID,
                this.props.bucketName,
                this.state.continuationToken,
                this.BULK_SIZE
            );

            this.setState({
                objects:            [...this.state.objects, ...resp.objects],
                continuationToken:  resp.continuation_token,
                hasMore:            !!(resp.continuation_token && resp.objects.length >= this.BULK_SIZE),
                loading:            false
            });

            this.props.onError(null);
        }
        catch ( err )
        {
            console.error(`Error loading objects: ${(err as Error).message}`, err);

            this.props.onError(err as Error);

            this.setState({
                loading:    false,
                hasMore:    false
            });
        }
    }


    private async _downloadObject ( object: ObjectV1 )
    {
        try
        {
            await this._downloadManagerService.download(
                this.props.profileID,
                this.props.bucketName,
                object.key
            );

            this.props.onError(null);
        }
        catch ( err )
        {
            console.error(`Error downloading object ${object.key}: ${(err as Error).message}`, err);

            this.props.onError(err as Error);
        }
    }


    public async componentDidMount ( ): Promise<void>
    {
        await this._load();

        this._s3ProfileService.updated().subscribe(this, this._load);
        this._s3ObjectService.updated().subscribe(this, this._load);
    }


    public async componentDidUpdate ( prevProps: ObjectListProps ): Promise<void>
    {
        if ( prevProps.bucketName === this.props.bucketName &&
             prevProps.profileID === this.props.profileID )
        {
            return;
        }

        await this._load();
    }


    public componentWillUnmount ( ): void
    {
        this._s3ProfileService.updated().unsubscribe(this);
        this._s3ObjectService.updated().unsubscribe(this);
    }


    public render ( )
    {
        return (
            <div className='ObjectList'>
                <div className='ObjectList-list'>
                    <div className='ObjectList-header'>
                        <div className='ObjectList-header-key'>Key</div>
                        <div className='ObjectList-header-owner'>Owner</div>
                        <div className='ObjectList-header-last-modified'>Last modified</div>
                        <div className='ObjectList-header-size'>Size</div>
                        <div className='ObjectList-header-actions'>
                            <FontAwesomeIcon icon={faSync} onClick={this._load} />
                        </div>
                    </div>

                    <div className='ObjectList-objects'>
                        <InfiniteScroll
                            pageStart={0}
                            loadMore={this._loadMore}
                            initialLoad={false}
                            hasMore={this.state.hasMore}
                            threshold={50}
                            useWindow={false}>
                            {this.props.bucketName && this.state.objects.length > 0 ?
                                this.state.objects.map( ( object ) => 
                                    <div
                                        className='ObjectList-object'
                                        key={object.key}>
                                        <div
                                            className='ObjectList-object-key'
                                            title={object.key}
                                            onClick={ ( ) => this.props.onSelectObject(object) }>
                                            {object.key}
                                        </div>

                                        <div
                                            className='ObjectList-object-owner'
                                            title={`${object.owner_name} (${object.owner_id})`}
                                            onClick={ ( ) => this.props.onSelectObject(object) }>
                                            {object.owner_name}
                                        </div>

                                        <div
                                            className='ObjectList-object-last-modified'
                                            title={object.last_modified}
                                            onClick={ ( ) => this.props.onSelectObject(object) }>
                                            {Moment(object.last_modified).format('YYYY-MM-DD HH:mm:ss')}
                                        </div>
            
                                        <div
                                            className='ObjectList-object-size'
                                            title={`${object.size} B`}
                                            onClick={ ( ) => this.props.onSelectObject(object) }>
                                            {Bytes.format(object.size, {unitSeparator: ' '})}
                                        </div>
                    
                                        <div
                                            className='ObjectList-object-actions'>
                                            <FontAwesomeIcon icon={faDownload} onClick={ ( ) => this._downloadObject(object) } />
                                            <FontAwesomeIcon icon={faTimes} onClick={ ( ) => this.props.onDeleteObject(object) } />
                                        </div>
                                    </div>
                                )
                            : null}

                            <Spinner active={this.state.loading} />
                            
                            {!this.state.loading && this.props.bucketName && this.state.objects.length === 0 ? 
                                <div className='ObjectList-empty'>No objects found.</div>
                            : null}

                            {!this.state.loading && !this.props.bucketName ? 
                                <div className='ObjectList-empty'>No bucket selected.</div>
                            : null}
                        </InfiniteScroll>
                    </div>
                </div>
            </div>
        );
    }
}
