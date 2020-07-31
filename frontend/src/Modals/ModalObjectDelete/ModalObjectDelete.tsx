import * as React from 'react';
import { Modal } from '../../Shared/Modal/Modal';
import { Button } from '../../Shared/Button/Button';
import { S3ObjectService, ObjectV1 } from '../../Shared/Service/ObjectService';

import './ModalObjectDelete.css';


export interface ModalObjectDeleteProps
{
    profileID:  number;
    bucketName: string;
    object:     ObjectV1;
    onClose:    ( ) => any;
    onSuccess:  ( ) => any;
    onError:    ( err: Error ) => any;
}


interface ModalObjectDeleteState
{
    loading:        boolean;
}


export class ModalObjectDelete extends React.Component<ModalObjectDeleteProps, ModalObjectDeleteState>
{
    private readonly _s3ObjectService: S3ObjectService;


    constructor ( props: ModalObjectDeleteProps )
    {
        super(props);

        this.state = {
            loading:    false
        };

        this._s3ObjectService  = S3ObjectService.getInstance();

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
            await this._s3ObjectService.deleteObject(
                this.props.profileID,
                this.props.bucketName,
                this.props.object.key
            );

            this.setState({
                loading:    false,
            });

            this.props.onSuccess();
        }
        catch ( err )
        {
            console.error(`Error deleting objects: ${err.message}`, err);
        
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
                title='Delete object'
                onClose={this.props.onClose}>
                <form onSubmit={this._onSubmit}>
                    <div className='ModalObjectDelete-text'>
                        Do you want to delete Object "{this.props.object.key}"?
                    </div>

                    <div className='ModalObjectDelete-actions'>
                        <Button type='submit'>
                            Delete
                        </Button>
                    </div>
                </form>
            </Modal>
        );
    }
}
