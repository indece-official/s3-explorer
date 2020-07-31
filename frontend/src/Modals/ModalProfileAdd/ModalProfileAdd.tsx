import * as React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPlus, faCaretDown, faCaretRight } from '@fortawesome/free-solid-svg-icons';
import { S3ProfileService, ProfileV1 } from '../../Shared/Service/ProfileService';
import { Modal } from '../../Shared/Modal/Modal';
import { Button } from '../../Shared/Button/Button';
import { InputText } from '../../Shared/Input/InputText';
import { InputCheckbox } from '../../Shared/Input/InputCheckbox';
import { Hint } from '../../Shared/Hint/Hint';

import './ModalProfileAdd.css';


export interface ModalProfileAddProps
{
    profile?:   ProfileV1;
    onClose:    ( ) => any;
    onSuccess:  ( profile: ProfileV1 ) => any;
    onError:    ( err: Error ) => any;
}



interface ModalProfileAddFormData
{
    name:       string;
    access_key: string;
    secret_key: string;
    region:     string;
    endpoint:   string;
    ssl:        boolean;
    path_style: boolean;
    buckets:    Array<string>;
}


interface ModalProfileAddFormErrors
{
    name?:       string;
    access_key?: string;
    secret_key?: string;
    region?:     string;
    endpoint?:   string;
    ssl?:        string;
    path_style?: string;
    buckets?:    Array<string>;
}


interface ModalProfileAddState
{
    editMode:       boolean;
    form:           ModalProfileAddFormData;
    formErrors:     ModalProfileAddFormErrors;
    advanced:       boolean;
    loading:        boolean;
}


export class ModalProfileAdd extends React.Component<ModalProfileAddProps, ModalProfileAddState>
{
    private readonly _s3ProfileService: S3ProfileService;


    constructor ( props: ModalProfileAddProps )
    {
        super(props);

        this.state = {
            editMode:       false,
            form: {
                name:       '',
                access_key: '',
                secret_key: '',
                region:     '',
                endpoint:   '',
                ssl:        true,
                path_style: false,
                buckets:    []
            },
            formErrors:     {},
            advanced:       false,
            loading:        false
        };

        this._s3ProfileService  = S3ProfileService.getInstance();

        this._onInputChange     = this._onInputChange.bind(this);
        this._onSubmit          = this._onSubmit.bind(this);
        this._addBucket         = this._addBucket.bind(this);
        this._toggleAdvanced    = this._toggleAdvanced.bind(this);
    }


    private _checkForm ( formData: ModalProfileAddFormData ): ModalProfileAddFormErrors | null
    {
        const formErrors: ModalProfileAddFormErrors = {};
        let hasErrors = false;

        if ( ! formData.name )
        {
            formErrors.name = 'Please enter a name';
            hasErrors = true;
        }

        if ( ! formData.access_key )
        {
            formErrors.access_key = 'Please enter the access key';
            hasErrors = true;
        }

        if ( ! formData.secret_key )
        {
            formErrors.secret_key = 'Please enter the secret key';
            hasErrors = true;
        }

        if ( hasErrors )
        {
            return formErrors;
        }

        return null;
    }


    private _onInputChange ( evt: any )
    {
        const target = evt.target;
        const value = target.type === 'checkbox' ? target.checked : target.value;
        const form = {
            ...this.state.form
        };

        const nameParts = target.name.split('.');
        let elem: any = form;
        for ( let i = 0; i < nameParts.length - 1; i++ )
        {
            elem = elem[nameParts[i]];
        }
        elem[nameParts[nameParts.length - 1]] = value;

        if ( this._checkForm(form) )
        {
            this.setState({
                form
            });
        }
        else
        {
            this.setState({
                form,
                formErrors: {}
            });
        }
    }


    private async _onSubmit ( evt: any ): Promise<void>
    {
        evt.preventDefault();

        const formErrors = this._checkForm(this.state.form);

        if ( formErrors )
        {
            this.setState({
                formErrors
            });

            return;
        }
        
        this.setState({
            formErrors: {},
            loading:    true
        });

        try
        {
            const data = {
                name:       this.state.form.name,
                access_key: this.state.form.access_key,
                secret_key: this.state.form.secret_key,
                region:     this.state.form.region,
                endpoint:   this.state.form.endpoint,
                ssl:        this.state.form.ssl,
                path_style: this.state.form.path_style,
                buckets:    this.state.form.buckets.filter( bucket => !!bucket )
            };

            
            let profileID = 0;
            
            if ( this.state.editMode )
            {
                await this._s3ProfileService.updateProfile(this.props.profile!.id, data);

                profileID = this.props.profile!.id;
            }
            else
            {
                profileID = await this._s3ProfileService.addProfile(data);
            }

            this.setState({
                loading:    false,
            });

            this.props.onSuccess({
                ...this.state.form,
                id:     profileID,
            });
        }
        catch ( err )
        {
            console.error(`Error adding profile: ${err.message}`, err);
        
            this.setState({
                loading:    false
            });

            this.props.onError(err);
        }
    }


    private _toggleAdvanced ( ): void
    {
        this.setState({
            advanced: !this.state.advanced
        })
    }


    private _addBucket ( ): void
    {
        const formData = this.state.form;

        formData.buckets.push('');

        this.setState({
            form: formData
        });
    }


    public async componentDidMount ( ): Promise<void>
    {
        if ( this.props.profile )
        {
            this.setState({
                editMode: true,
                form: {
                    name:       this.props.profile.name,
                    access_key: this.props.profile.access_key,
                    secret_key: this.props.profile.secret_key,
                    region:     this.props.profile.region,
                    endpoint:   this.props.profile.endpoint,
                    ssl:        this.props.profile.ssl,
                    path_style: this.props.profile.path_style,
                    buckets:    this.props.profile.buckets
                }
            });
        }
    }


    public render ( )
    {
        return (
            <Modal
                title={this.state.editMode ? 'Edit a profile' : 'Add a new profile'}
                onClose={this.props.onClose}>
                <form onSubmit={this._onSubmit}>    
                    <div className='ModalProfileAdd-inputs'>
                        <InputText
                            label='Name'
                            name='name'
                            error={this.state.formErrors.name}
                            value={this.state.form.name}
                            onChange={this._onInputChange}
                        />

                        <InputText
                            label='Access-Key'
                            name='access_key'
                            error={this.state.formErrors.access_key}
                            value={this.state.form.access_key}
                            onChange={this._onInputChange}
                        />

                        <InputText
                            label='Secret-Key'
                            name='secret_key'
                            error={this.state.formErrors.secret_key}
                            value={this.state.form.secret_key}
                            onChange={this._onInputChange}
                        />

                        <InputText
                            label='Region'
                            name='region'
                            error={this.state.formErrors.region}
                            value={this.state.form.region}
                            onChange={this._onInputChange}
                        />

                        <div
                            className='ModalProfileAdd-advanced-toggle'
                            onClick={this._toggleAdvanced}>
                            {this.state.advanced ?
                                <FontAwesomeIcon
                                    icon={faCaretDown}
                                    fixedWidth={true}
                                />
                            :
                                <FontAwesomeIcon
                                    icon={faCaretRight}
                                    fixedWidth={true}
                                />
                            }
                            &nbsp; Advanced options
                        </div>
                        
                        {this.state.advanced ?
                            <div className='ModalProfileAdd-advanced'>
                                <InputText
                                    label='Endpoint'
                                    name='endpoint'
                                    error={this.state.formErrors.endpoint}
                                    value={this.state.form.endpoint}
                                    onChange={this._onInputChange}
                                />

                                <InputCheckbox
                                    label='Use SSL'
                                    name='ssl'
                                    value={this.state.form.ssl}
                                    onChange={this._onInputChange}
                                />

                                <InputCheckbox
                                    label='Use path-style URIs'
                                    name='path_style'
                                    value={this.state.form.path_style}
                                    onChange={this._onInputChange}
                                />

                                <div className='ModalProfileAdd-buckets'>
                                    <div className='ModalProfileAdd-buckets-title'>
                                        Preconfigured Buckets:
                                        &emsp;
                                        <Hint
                                            text={'Buckets must only be preconfigured if the profile doesn\'t have the permission to list buckets.'}
                                        />
                                    </div>
                                    
                                    <div className='ModalProfileAdd-buckets-bucket'>
                                        {this.state.form.buckets.map( ( bucket, i ) => 
                                            <InputText
                                                key={i}
                                                label='Bucket'
                                                name={`buckets.${i}`}
                                                error={this.state.formErrors.buckets && this.state.formErrors.buckets[i]}
                                                value={this.state.form.buckets[i]}
                                                onChange={this._onInputChange}
                                            />
                                        )}
                                    </div>

                                    <div className='ModalProfileAdd-buckets-add'>
                                        <FontAwesomeIcon
                                            icon={faPlus}
                                            onClick={this._addBucket}
                                            title='Add a preconfigured bucket'
                                        />
                                    </div>
                                </div>
                            </div>
                        : null}
                    </div>

                    <div className='ModalProfileAdd-actions'>
                        <Button type='submit'>
                            Save
                        </Button>
                    </div>
                </form>
            </Modal>
        );
    }
}
