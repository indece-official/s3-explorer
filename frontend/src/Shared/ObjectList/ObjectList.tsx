import * as React from 'react';
import Bytes from 'bytes';
import Moment from 'moment';
import { ObjectV1, S3ObjectService } from '../Service/ObjectService';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSync, faTimes, faDownload } from '@fortawesome/free-solid-svg-icons';

import './ObjectList.css';
import { DownloadManagerService } from '../DownloadManager/DownloadManagerService';


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
    objects: Array<ObjectV1>;
}


export class ObjectList extends React.Component<ObjectListProps, ObjectListState>
{
    private readonly _s3ObjectService:          S3ObjectService;
    private readonly _downloadManagerService:   DownloadManagerService;


    constructor ( props: ObjectListProps )
    {
        super(props);

        this.state = {
            objects: []
        };

        this._s3ObjectService = S3ObjectService.getInstance();
        this._downloadManagerService = DownloadManagerService.getInstance();

        this._load = this._load.bind(this);
    }


    private async _load ( ): Promise<void>
    {
        if ( !this.props.profileID || !this.props.bucketName )
        {
            this.setState({
                objects: []
            });

            return;
        }

        try
        {
            const objects = await this._s3ObjectService.getObjects(this.props.profileID, this.props.bucketName);

            this.setState({
                objects
            });

            this.props.onError(null);
        }
        catch ( err )
        {
            console.error(`Error loading objects: ${err.message}`, err);

            this.props.onError(err);
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
            console.error(`Error downloading object ${object.key}: ${err.message}`, err);

            this.props.onError(err);
        }
    }


    public async componentDidMount ( ): Promise<void>
    {
        await this._load();

        this._s3ObjectService.updated().subscribe(this, this._load.bind(this));
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
                        : (this.props.bucketName && this.state.objects.length == 0 ? 
                            <div className='ObjectList-empty'>No objects found.</div>
                        : 
                            <div className='ObjectList-empty'>No bucket selected.</div>
                        )}
                    </div>
                </div>
            </div>
        );
    }
}
