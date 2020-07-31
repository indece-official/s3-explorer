import * as React from 'react';
import { Modal } from '../../Shared/Modal/Modal';
import { ProfileV1, S3ProfileService } from '../../Shared/Service/ProfileService';
import { Button } from '../../Shared/Button/Button';

import './ModalProfileSelect.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faTimes, faPen } from '@fortawesome/free-solid-svg-icons';


export interface ModalProfileSelectProps
{
    onSelectProfile:    ( profile: ProfileV1 ) => any;
    onDeleteProfile:    ( profile: ProfileV1 ) => any;
    onEditProfile:      ( profile: ProfileV1 ) => any;
    onAddProfile:       ( ) => any;
    onClose:            ( ) => any;
    onError:            ( err: Error ) => any;
}


interface ModalProfileSelectState
{
    profiles: Array<ProfileV1>;
}


export class ModalProfileSelect extends React.Component<ModalProfileSelectProps, ModalProfileSelectState>
{
    private readonly _s3ProfileService: S3ProfileService;


    constructor ( props: ModalProfileSelectProps )
    {
        super(props);

        this.state = {
            profiles: []
        };

        this._s3ProfileService = S3ProfileService.getInstance();
    }


    private async _load ( ): Promise<void>
    {
        try
        {
            const profiles = await this._s3ProfileService.getProfiles();

            this.setState({
                profiles
            });
        }
        catch ( err )
        {
            console.error(`Error loading profiles: ${err.message}`, err);

            this.props.onError(err);
        }
    }


    public async componentDidMount ( ): Promise<void>
    {
       await this._load();

       this._s3ProfileService.updated().subscribe(this, this._load.bind(this));
    }


    public componentWillUnmount ( ): void
    {
        this._s3ProfileService.updated().unsubscribe(this);
    }


    public render ( )
    {
        return (
            <Modal
                title='Select a profile'
                onClose={this.props.onClose}>
                <div className='ModalProfileSelect-actions'>
                    <Button
                        onClick={this.props.onAddProfile}
                        type='button'>
                        Add profile
                    </Button>
                </div> 
 
                <div className='ModalProfileSelect-profiles'>
                    {this.state.profiles.map( ( profile ) => 
                        <div
                            key={profile.id}
                            className='ModalProfileSelect-profile'>
                            <div
                                className='ModalProfileSelect-profile-name'
                                onClick={ ( ) => this.props.onSelectProfile(profile) }>
                                {profile.name}
                            </div>

                            <div
                                className='ModalProfileSelect-profile-actions'>
                                <FontAwesomeIcon
                                    icon={faPen}
                                    title='Edit profile'
                                    onClick={ ( ) => this.props.onEditProfile(profile) }
                                />
                                
                                <FontAwesomeIcon
                                    icon={faTimes}
                                    title='Delete profile'
                                    onClick={ ( ) => this.props.onDeleteProfile(profile) }
                                />
                            </div>
                        </div>
                    )}
                </div>
            </Modal>
        )
    }
}