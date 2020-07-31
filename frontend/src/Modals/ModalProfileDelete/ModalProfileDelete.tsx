import * as React from 'react';
import { Modal } from '../../Shared/Modal/Modal';
import { Button } from '../../Shared/Button/Button';
import { ProfileV1,
         S3ProfileService } from '../../Shared/Service/ProfileService';

import './ModalProfileDelete.css';


export interface ModalProfileDeleteProps
{
    profile:    ProfileV1;
    onClose:    ( ) => any;
    onSuccess:  ( ) => any;
    onError:    ( err: Error ) => any;
}


interface ModalProfileDeleteState
{
    loading:        boolean;
}


export class ModalProfileDelete extends React.Component<ModalProfileDeleteProps, ModalProfileDeleteState>
{
    private readonly _s3ProfileService: S3ProfileService;


    constructor ( props: ModalProfileDeleteProps )
    {
        super(props);

        this.state = {
            loading:    false
        };

        this._s3ProfileService   = S3ProfileService.getInstance();

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
            await this._s3ProfileService.deleteProfile(
                this.props.profile.id
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
                title='Delete profile'
                onClose={this.props.onClose}>
                <form onSubmit={this._onSubmit}>
                    <div className='ModalProfileDelete-text'>
                        Do you want to delete Profile "{this.props.profile.name}"?<br />
                        <br />
                        This won't delete the buckets and objects accessable via this profile.
                    </div>

                    <div className='ModalProfileDelete-actions'>
                        <Button type='submit'>
                            Delete
                        </Button>
                    </div>
                </form>
            </Modal>
        );
    }
}
