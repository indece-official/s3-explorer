import * as React from 'react';
import { Modal } from '../../Shared/Modal/Modal';
import { Button } from '../../Shared/Button/Button';
import { BucketV1,
         S3BucketService } from '../../Shared/Service/BucketService';

import './ModalBucketDelete.css';


export interface ModalBucketDeleteProps
{
    profileID:  number;
    bucket:     BucketV1;
    onClose:    ( ) => any;
    onSuccess:  ( ) => any;
    onError:    ( err: Error ) => any;
}


interface ModalBucketDeleteState
{
    loading:        boolean;
}


export class ModalBucketDelete extends React.Component<ModalBucketDeleteProps, ModalBucketDeleteState>
{
    private readonly _s3BucketService: S3BucketService;


    constructor ( props: ModalBucketDeleteProps )
    {
        super(props);

        this.state = {
            loading:    false
        };

        this._s3BucketService   = S3BucketService.getInstance();

        this._onSubmit          = this._onSubmit.bind(this);
    }


    private async _onSubmit ( evt: any ): Promise<void>
    {
        evt.preventDefault();

        this.setState({
            loading:    true
        });

        try
        {
            await this._s3BucketService.deleteBucket(
                this.props.profileID,
                this.props.bucket.name
            );

            this.setState({
                loading:    false,
            });

            this.props.onSuccess();
        }
        catch ( err )
        {
            console.error(`Error deleting bucket: ${err.message}`, err);
        
            this.setState({
                loading:    false
            });

            this.props.onError(err);
        }
    }


    public render ( )
    {
        return (
            <Modal
                title='Delete bucket'
                onClose={this.props.onClose}>
                <form onSubmit={this._onSubmit}>
                    <div className='ModalBucketDelete-text'>
                        Do you want to delete Bucket "{this.props.bucket.name}"?<br />
                        <br />
                        The bucket must be empty before you can delete it.
                    </div>

                    <div className='ModalBucketDelete-actions'>
                        <Button type='submit'>
                            Delete
                        </Button>
                    </div>
                </form>
            </Modal>
        );
    }
}
