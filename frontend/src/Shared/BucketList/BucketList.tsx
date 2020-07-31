import * as React from 'react';
import { S3BucketService, BucketV1 } from '../Service/BucketService';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSync, faPlus, faTimes } from '@fortawesome/free-solid-svg-icons';

import './BucketList.css';


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
}


export class BucketList extends React.Component<BucketListProps, BucketListState>
{
    private readonly _s3BucketService: S3BucketService;


    constructor ( props: BucketListProps )
    {
        super(props);

        this.state = {
            buckets:    []
        };

        this._s3BucketService = S3BucketService.getInstance();
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

    
    private _onSelectBucket ( bucket: BucketV1 ): void
    {
        this.props.onSelectBucket(bucket);
    }


    public async componentDidMount ( ): Promise<void>
    {
        await this._load();

        this._s3BucketService.updated().subscribe(this, this._load.bind(this));
    }
    

    public async componentDidUpdate ( prevProps: BucketListProps ): Promise<void>
    {
        if ( prevProps.profileID === this.props.profileID )
        {
            return;
        }

        await this._load();
    }


    public componentWillUnmount ( ): void
    {
        this._s3BucketService.updated().unsubscribe(this);
    }


    public render ( )
    {
        return (
            <div className='BucketList'>
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
                    : (this.props.profileID && this.state.buckets.length == 0 ? 
                        <div className='BucketList-empty'>No buckets found.</div>
                    : 
                        <div className='BucketList-empty'>No profile selected.</div>
                    )}
                </div>
            </div>
        );
    }
}
