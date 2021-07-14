import * as React from 'react';
import Bytes from 'bytes';
import { S3BucketService, BucketV1, BucketStatsV1 } from '../Service/BucketService';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSync, faPlus, faTimes } from '@fortawesome/free-solid-svg-icons';
import { Button } from '../Button/Button';

import './BucketList.css';
import { MiniSpinner } from '../MiniSpinner/MiniSpinner';
import { S3ObjectService } from '../Service/ObjectService';


export interface BucketListProps
{
    profileID:      number;
    selectedBucket: string;
    onAddBucket:    ( ) => any;
    onSelectBucket: ( bucket: BucketV1 ) => any;
    onDeleteBucket: ( bucket: BucketV1 ) => any;
    onError:        ( err: Error | null ) => any;
}


interface BucketListState
{
    buckets:        Array<BucketV1>;
    bucketStats:    BucketStatsV1 | null;
}


export class BucketList extends React.Component<BucketListProps, BucketListState>
{
    private readonly _s3BucketService: S3BucketService;
    private readonly _s3ObjectService: S3ObjectService;
    private _latestBucketStatsQueryID: number;


    constructor ( props: BucketListProps )
    {
        super(props);

        this.state = {
            buckets:        [],
            bucketStats:    null
        };

        this._latestBucketStatsQueryID = 1;
        this._s3BucketService = S3BucketService.getInstance();
        this._s3ObjectService = S3ObjectService.getInstance();
    }


    private async _load ( ): Promise<void>
    {
        if ( ! this.props.profileID )
        {
            this.setState({
                buckets: []
            });

            return;
        }

        try
        {
            const buckets = await this._s3BucketService.getBuckets(this.props.profileID);

            this.setState({
                buckets
            });
            
            this.props.onError(null);
        }
        catch ( err )
        {
            console.error(`Error loading buckets: ${err.message}`, err);

            this.props.onError(err);
        }
    }
   
   
    private async _loadBucketStats ( force: boolean ): Promise<void>
    {
        if ( !this.props.profileID || !this.props.selectedBucket )
        {
            this.setState({
                bucketStats: null
            });

            return;
        }

        try
        {
            this.setState({
                bucketStats: null
            });

            this._latestBucketStatsQueryID++;

            const queryID = this._latestBucketStatsQueryID;

            const bucketStats = await this._s3BucketService.getBucketStats(this.props.profileID, this.props.selectedBucket, force);

            if ( queryID < this._latestBucketStatsQueryID )
            {
                // Other query was fired in between, cancel

                return;
            }

            this.setState({
                bucketStats
            });
            
            this.props.onError(null);
        }
        catch ( err )
        {
            console.error(`Error loading bucket stats: ${err.message}`, err);

            this.props.onError(err);
        }
    }

    
    private _onSelectBucket ( bucket: BucketV1 ): void
    {
        this.props.onSelectBucket(bucket);
    }


    public async componentDidMount ( ): Promise<void>
    {
        await this._load();
        await this._loadBucketStats(false);

        this._s3BucketService.updated().subscribe(this, this._load.bind(this));
        this._s3ObjectService.updated().subscribe(this, ( ) => this._loadBucketStats(false) );
    }
    

    public async componentDidUpdate ( prevProps: BucketListProps ): Promise<void>
    {
        if ( prevProps.profileID !== this.props.profileID )
        {
            await this._load();
        }
        
        if ( prevProps.selectedBucket !== this.props.selectedBucket )
        {
            await this._loadBucketStats(false);
        }
    }


    public componentWillUnmount ( ): void
    {
        this._s3BucketService.updated().unsubscribe(this);
        this._s3ObjectService.updated().unsubscribe(this);
    }


    public render ( )
    {
        return (
            <div className='BucketList'>
                <div className='BucketList-list'>
                    <div className='BucketList-header'>
                        <div className='BucketList-header-title'>
                            Buckets
                        </div>

                        <div
                            className='BucketList-header-actions'>
                            <FontAwesomeIcon icon={faPlus} onClick={this.props.onAddBucket} />
                            <FontAwesomeIcon icon={faSync} onClick={this._load} />
                        </div>
                    </div>

                    <div className='BucketList-buckets'>
                        {this.props.profileID && this.state.buckets.length > 0 ? 
                            this.state.buckets.map( ( bucket ) => 
                                <div
                                    className={'BucketList-bucket' + (bucket.name === this.props.selectedBucket ? ' BucketList-bucket-selected' : '')}
                                    key={bucket.name}>
                                    <div
                                        className='BucketList-bucket-name'
                                        onClick={ ( ) => this._onSelectBucket(bucket) }>
                                        {bucket.name}
                                    </div>
                                    
                                    <div
                                        className='BucketList-bucket-actions'
                                        onClick={ ( ) => this.props.onDeleteBucket(bucket) }>
                                        <FontAwesomeIcon icon={faTimes} title='Delete bucket' />
                                    </div>
                                </div>
                            )
                        : (this.props.profileID && this.state.buckets.length === 0 ? 
                            <div className='BucketList-empty'>No buckets found.</div>
                        : 
                            <div className='BucketList-empty'>No profile selected.</div>
                        )}
                    </div>
                </div>

                {this.props.selectedBucket ?
                    <div className='BucketList-currbucket'>
                        <div className='BucketList-currbucket-header'>
                            <div className='BucketList-currbucket-name'>
                                {this.props.selectedBucket}
                            </div>
                        </div>

                        <div className='BucketList-currbucket-stats'>
                            {this.state.bucketStats && this.state.bucketStats.complete ?
                                <>
                                    <div className='BucketList-currbucket-stat'>Total files: {this.state.bucketStats.files}</div>
                                    <div className='BucketList-currbucket-stat'>Total size: {Bytes(this.state.bucketStats.size)}</div>
                                </>
                            : (this.state.bucketStats && !this.state.bucketStats.complete) ?
                                <>
                                    <div className='BucketList-currbucket-stat'>Total files: &gt; {this.state.bucketStats.files}</div>
                                    <div className='BucketList-currbucket-stat'>Total size: &gt; {Bytes(this.state.bucketStats.size)}</div>
                                    <Button onClick={ ( ) => this._loadBucketStats(true) }>Count total</Button>
                                </>
                            :
                                <>
                                    <div className='BucketList-currbucket-stat'>Total files: <MiniSpinner /></div>
                                    <div className='BucketList-currbucket-stat'>Total size: <MiniSpinner /></div>
                                </>
                            }
                        </div>
                    </div>
                : null}
            </div>
        );
    }
}
