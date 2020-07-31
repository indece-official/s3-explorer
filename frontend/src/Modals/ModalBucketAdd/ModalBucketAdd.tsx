import * as React from 'react';
import { S3BucketService, BucketV1 } from '../../Shared/Service/BucketService';
import { Modal } from '../../Shared/Modal/Modal';
import { Button } from '../../Shared/Button/Button';
import { InputText } from '../../Shared/Input/InputText';

import './ModalBucketAdd.css';


export interface ModalBucketAddProps
{
    profileID:  number;
    onClose:    ( ) => any;
    onSuccess:  ( bucket: BucketV1 ) => any;
    onError:    ( err: Error ) => any;
}



interface ModalBucketAddFormData
{
    name:       string;
}


interface ModalBucketAddFormErrors
{
    name?:       string;
}


interface ModalBucketAddState
{
    form:           ModalBucketAddFormData;
    formErrors:     ModalBucketAddFormErrors;
    loading:        boolean;
}


export class ModalBucketAdd extends React.Component<ModalBucketAddProps, ModalBucketAddState>
{
    private readonly _s3BucketService: S3BucketService;


    constructor ( props: ModalBucketAddProps )
    {
        super(props);

        this.state = {
            form: {
                name:       ''
            },
            formErrors:     {},
            loading:        false
        };

        this._s3BucketService  = S3BucketService.getInstance();

        this._onInputChange     = this._onInputChange.bind(this);
        this._onSubmit          = this._onSubmit.bind(this);
    }


    private _checkForm ( formData: ModalBucketAddFormData ): ModalBucketAddFormErrors | null
    {
        const formErrors: ModalBucketAddFormErrors = {};
        let hasErrors = false;

        if ( ! formData.name )
        {
            formErrors.name = 'Please enter a name';
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
        const name = target.name;
        const form = {
            ...this.state.form,
            [name]: value
        };

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

        if ( this.state.loading )
        {
            return;
        }

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
            await this._s3BucketService.addBucket(this.props.profileID, {
                name:       this.state.form.name
            });

            this.setState({
                loading:    false,
            });

            this.props.onSuccess({
                name: this.state.form.name
            });
        }
        catch ( err )
        {
            console.error(`Error adding bucket: ${err.message}`, err);
        
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
                title='Add a new bucket'
                onClose={this.props.onClose}>
                <form onSubmit={this._onSubmit}>
                    <div className='ModalBucketAdd-inputs'>
                        <InputText
                            label='Name'
                            name='name'
                            error={this.state.formErrors.name}
                            value={this.state.form.name}
                            onChange={this._onInputChange}
                        />
                    </div>
                    
                    <div className='ModalBucketAdd-actions'>
                        <Button
                            type='submit'
                            disabled={this.state.loading}>
                            Save
                        </Button>
                    </div>
                </form>
            </Modal>
        );
    }
}
