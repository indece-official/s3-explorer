import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faAngleRight } from '@fortawesome/free-solid-svg-icons';
import { NavBar } from './Shared/NavBar/NavBar';
import { BucketList } from './Shared/BucketList/BucketList';
import { ObjectList } from './Shared/ObjectList/ObjectList';
import { ModalObjectDetails } from './Modals/ModalObjectDetails/ModalObjectDetails';
import { ModalProfileSelect } from './Modals/ModalProfileSelect/ModalProfileSelect';
import { ModalProfileAdd } from './Modals/ModalProfileAdd/ModalProfileAdd';
import { ModalBucketAdd } from './Modals/ModalBucketAdd/ModalBucketAdd';
import { ModalObjectAdd } from './Modals/ModalObjectAdd/ModalObjectAdd';
import { S3ProfileService, ProfileV1 } from './Shared/Service/ProfileService';
import { ModalObjectDelete } from './Modals/ModalObjectDelete/ModalObjectDelete';
import { ObjectV1 } from './Shared/Service/ObjectService';
import { DownloadManager } from './Shared/DownloadManager/DownloadManager';
import { ModalAbout } from './Modals/ModalAbout/ModalAbout';
import { BucketV1 } from './Shared/Service/BucketService';
import { ModalBucketDelete } from './Modals/ModalBucketDelete/ModalBucketDelete';
import { ModalProfileDelete } from './Modals/ModalProfileDelete/ModalProfileDelete';

import './App.css';


interface AppState
{
    error:              Error | null;
    selectedBucket:     BucketV1 | null;
    selectedProfile:    ProfileV1 | null;
    showSelectProfile:  boolean;
    showAddProfile:     boolean;
    showAddBucket:      boolean;
    showAddObject:      boolean;
    showAbout:          boolean;
    showDetailsObject:  ObjectV1 | null;
    showDeleteObject:   ObjectV1 | null;
    showDeleteBucket:   BucketV1 | null;
    showEditProfile:    ProfileV1 | null;
    showDeleteProfile:  ProfileV1 | null;
}


export class App extends React.Component<{}, AppState>
{
    private readonly _s3ProfileService: S3ProfileService;



    constructor ( props: any )
    {
        super(props);

        this.state = {
            error:              null,
            selectedBucket:     null,
            selectedProfile:    null,
            showSelectProfile:  false,
            showAddProfile:     false,
            showAddBucket:      false,
            showAddObject:      false,
            showAbout:          false,
            showDetailsObject:  null,
            showDeleteObject:   null,
            showDeleteBucket:   null,
            showEditProfile:    null,
            showDeleteProfile:  null
        };

        this._s3ProfileService = S3ProfileService.getInstance();

        this.onError            = this.onError.bind(this);
        this.onSelectProfile    = this.onSelectProfile.bind(this);
        this.onSelectBucket     = this.onSelectBucket.bind(this);
        this.showDetailsObject  = this.showDetailsObject.bind(this);
        this.showDeleteObject   = this.showDeleteObject.bind(this);
        this.showDeleteBucket   = this.showDeleteBucket.bind(this);
        this.showEditProfile    = this.showEditProfile.bind(this);
        this.showDeleteProfile  = this.showDeleteProfile.bind(this);
    }


    private onError ( error: Error | null )
    {
        this.setState({
            error
        });
    }


    private onSelectProfile ( profile: ProfileV1 )
    {
        this.setState({
            selectedProfile:    profile,
            selectedBucket:     null,
            showSelectProfile:  false,
            showAddProfile:     false,
            showDetailsObject:  null,
            showDeleteObject:   null,
            showEditProfile:    null,
            showDeleteProfile:  null
        });
    }
    
    
    private onSelectBucket ( bucket: BucketV1 )
    {
        this.setState({
            selectedBucket:     bucket,
            showAddBucket:      false,
            showDetailsObject:  null,
            showDeleteObject:   null
        });
    }
    
    
    private showDetailsObject ( object: ObjectV1 | null )
    {
        this.setState({
            showDetailsObject:  object,
            showDeleteObject:   null
        });
    }
    
    
    private showDeleteObject ( object: ObjectV1 | null )
    {
        this.setState({
            showDetailsObject:  null,
            showDeleteObject:   object
        });
    }


    private showSelectProfile ( visible: boolean )
    {
        this.setState({
            showSelectProfile: visible
        });
    }


    private showAddProfile ( visible: boolean )
    {
        this.setState({
            showAddProfile:     visible,
            showSelectProfile:  false
        });
    }


    private showDeleteProfile ( profile: ProfileV1 | null )
    {
        this.setState({
            showDeleteProfile:  profile
        });
    }


    private showEditProfile ( profile: ProfileV1 | null )
    {
        this.setState({
            showEditProfile:  profile
        });
    }


    private showAddBucket ( visible: boolean )
    {
        this.setState({
            showAddBucket:  visible
        });
    }


    private showDeleteBucket ( bucket: BucketV1 | null, deleted?: boolean )
    {
        if ( deleted &&
             this.state.showDeleteBucket &&
             this.state.selectedBucket &&
             this.state.selectedBucket.name === this.state.showDeleteBucket.name )
        {
            this.setState({
                showAddProfile:     false,
                showDeleteBucket:   bucket,
                selectedBucket:     null
            });

            return;
        }

        this.setState({
            showAddProfile:     false,
            showDeleteBucket:   bucket
        });
    }
    
    
    private showAddObject ( visible: boolean )
    {
        this.setState({
            showAddObject:  visible
        });
    }
    
   
    private showAbout ( visible: boolean )
    {
        this.setState({
            showAbout:  visible
        });
    }
    

    public async componentDidMount ( )
    {
        try
        {
            const profiles = await this._s3ProfileService.getProfiles();

            if ( profiles.length === 0 )
            {
                this.setState({
                    showAddProfile: true
                });

                return;
            }

            this.setState({
                selectedProfile: profiles[0]
            });
        }
        catch ( err )
        {
            console.error(`can't load profiles: ${err.message}`, err);

            this.setState({
                error: new Error(`Can't load profiles`)
            });
        }
    }


    public render ( )
    {
        return (
            <div className="App">
                <NavBar
                    onSelectProfile={ ( ) => this.showSelectProfile(true) }
                    onAbout={ ( ) => this.showAbout(true) }
                    onAddObject={ ( ) => this.showAddObject(true) }
                />

                <DownloadManager
                    onError={this.onError}
                />

                {this.state.error ?
                    <div className='App-error'>
                        Error: {this.state.error.message}
                    </div>
                : null}

                <div className='App-path'>
                    Path: 
                    {this.state.selectedProfile ?
                        <span>
                            &nbsp;
                            {this.state.selectedProfile.name}
                            &nbsp;
                            <FontAwesomeIcon icon={faAngleRight} />
                        </span>
                    : null}
                    {this.state.selectedBucket ?
                        <span>
                            &nbsp;
                            {this.state.selectedBucket.name}
                            &nbsp;
                            <FontAwesomeIcon icon={faAngleRight} />
                        </span>
                    : null}
                </div>

                <div className='App-content'>
                    <div className='App-content-lists'>
                        <BucketList
                            profileID={this.state.selectedProfile ? this.state.selectedProfile.id : 0}
                            selectedBucket={this.state.selectedBucket ? this.state.selectedBucket.name : ''}
                            onSelectBucket={this.onSelectBucket}
                            onAddBucket={ ( ) => this.showAddBucket(true) }
                            onDeleteBucket={this.showDeleteBucket}
                            onError={this.onError}
                        />

                        <ObjectList
                            profileID={this.state.selectedProfile ? this.state.selectedProfile.id : 0}
                            bucketName={this.state.selectedBucket ? this.state.selectedBucket.name : ''}
                            onSelectObject={this.showDetailsObject}
                            onDeleteObject={this.showDeleteObject}
                            onError={this.onError}
                        />
                    </div>
                </div>

                {this.state.showDetailsObject ? 
                    <ModalObjectDetails
                        profileID={this.state.selectedProfile ? this.state.selectedProfile.id : 0}
                        bucketName={this.state.selectedBucket ? this.state.selectedBucket.name : ''}
                        object={this.state.showDetailsObject}
                        onClose={ ( ) => this.showDetailsObject(null) }
                        onError={this.onError}
                    />
                : null}
                
                {this.state.showDeleteObject ? 
                    <ModalObjectDelete
                        profileID={this.state.selectedProfile ? this.state.selectedProfile.id : 0}
                        bucketName={this.state.selectedBucket ? this.state.selectedBucket.name : ''}
                        object={this.state.showDeleteObject}
                        onClose={ ( ) => this.showDeleteObject(null) }
                        onSuccess={ ( ) => this.showDeleteObject(null) }
                        onError={this.onError}
                    />
                : null}
                
                {this.state.showSelectProfile ? 
                    <ModalProfileSelect
                        onClose={ ( ) => this.showSelectProfile(false) }
                        onSelectProfile={this.onSelectProfile}
                        onEditProfile={this.showEditProfile}
                        onDeleteProfile={this.showDeleteProfile}
                        onAddProfile={ ( ) => this.showAddProfile(true) }
                        onError={this.onError}
                    />
                : null}

                {this.state.showAddProfile ? 
                    <ModalProfileAdd
                        onClose={ ( ) => this.showAddProfile(false) }
                        onSuccess={this.onSelectProfile}
                        onError={this.onError}
                    />
                : null}
                
                {this.state.showEditProfile ? 
                    <ModalProfileAdd
                        profile={this.state.showEditProfile}
                        onClose={ ( ) => this.showEditProfile(null) }
                        onSuccess={this.onSelectProfile}
                        onError={this.onError}
                    />
                : null}

                {this.state.showDeleteProfile ? 
                    <ModalProfileDelete
                        profile={this.state.showDeleteProfile}
                        onClose={ ( ) => this.showDeleteProfile(null) }
                        onSuccess={ ( ) => this.showDeleteProfile(null) }
                        onError={this.onError}
                    />
                : null}

                {this.state.showAddBucket ? 
                    <ModalBucketAdd
                        profileID={this.state.selectedProfile ? this.state.selectedProfile.id : 0}
                        onClose={ ( ) => this.showAddBucket(false) }
                        onSuccess={this.onSelectBucket}
                        onError={this.onError}
                    />
                : null}

                {this.state.showAddObject ? 
                    <ModalObjectAdd
                        profileID={this.state.selectedProfile ? this.state.selectedProfile.id : 0}
                        bucketName={this.state.selectedBucket ? this.state.selectedBucket.name : ''}
                        onClose={ ( ) => this.showAddObject(false) }
                        onSuccess={ ( ) => this.showAddObject(false) }
                        onError={this.onError}
                    />
                : null}

                {this.state.showDeleteBucket ? 
                    <ModalBucketDelete
                        profileID={this.state.selectedProfile ? this.state.selectedProfile.id : 0}
                        bucket={this.state.showDeleteBucket}
                        onClose={ ( ) => this.showDeleteBucket(null, false) }
                        onSuccess={ ( ) => this.showDeleteBucket(null, true) }
                        onError={this.onError}
                    />
                : null}
                
                {this.state.showAbout ? 
                    <ModalAbout
                        onClose={ ( ) => this.showAbout(false) }
                    />
                : null}
            </div>
        );
    }
}

export default App;
