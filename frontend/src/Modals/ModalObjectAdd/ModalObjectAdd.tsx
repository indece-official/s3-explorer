import * as React from 'react';
import { Modal } from '../../Shared/Modal/Modal';
import { Button } from '../../Shared/Button/Button';
import { S3ObjectService } from '../../Shared/Service/ObjectService';
import { InputFile } from '../../Shared/Input/InputFile';
import { ProgressBar } from '../../Shared/ProgressBar/ProgressBar';

import './ModalObjectAdd.css';
import { InputText } from '../../Shared/Input/InputText';


export interface ModalObjectAddProps
{
    profileID:  number;
    bucketName: string;
    onClose:    ( ) => any;
    onSuccess:  ( ) => any;
    onError:    ( err: Error | null ) => any;
}



interface ModalObjectAddFormData
{
    files:      FileList | null;
    filename:   string;
}


interface ModalObjectAddFormErrors
{
    files?:     string;
    filename?:  string;
}


interface ModalObjectAddState
{
    form:           ModalObjectAddFormData;
    formErrors:     ModalObjectAddFormErrors;
    loading:        boolean;
    progress:       ProgressEvent | null;
}


export class ModalObjectAdd extends React.Component<ModalObjectAddProps, ModalObjectAddState>
{
    private readonly _s3ObjectService: S3ObjectService;


    constructor ( props: ModalObjectAddProps )
    {
        super(props);

        this.state = {
            form: {
                files:      null,
                filename:   ''
            },
            formErrors:     {},
            loading:        false,
            progress:       null
        };

        this._s3ObjectService  = S3ObjectService.getInstance();

        this._onInputChange     = this._onInputChange.bind(this);
        this._onSubmit          = this._onSubmit.bind(this);
        this._onProgress        = this._onProgress.bind(this);
    }


    private _checkForm ( formData: ModalObjectAddFormData ): ModalObjectAddFormErrors | null
    {
        const formErrors: ModalObjectAddFormErrors = {};
        let hasErrors = false;

        if ( ! formData.files )
        {
            formErrors.files = 'Please select a file';
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
        let value = null;
        
        switch ( target.type )
        {
            case 'checkbox':
                value = target.checked;
                break;
            case 'file':
                value = target.files;
                break;
            default:
                value = target.value;
                break;
        }

        const name = target.name;
        const form = {
            ...this.state.form,
            [name]: value
        };

        console.log(name, value);

        if ( name === 'files' )
        {
            form.filename = value.length === 1 ? value[0].name : '';
        }

        console.log(form);

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


    private _onProgress ( evt: ProgressEvent ): void
    {
        this.setState({
            progress: evt
        });
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
            for ( let i = 0; i < this.state.form.files!.length; i++ )
            {
                const file = this.state.form.files!.item(i);
                if ( ! file )
                {
                    continue;
                }

                await this._s3ObjectService.addObject(
                    this.props.profileID,
                    this.props.bucketName,
                    file,
                    this.state.form.filename || null,
                    this._onProgress
                );
            }

            this.setState({
                loading:    false,
            });

            this.props.onSuccess();
            this.props.onError(null);
        }
        catch ( err )
        {
            console.error(`Error adding objects: ${err.message}`, err);
        
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
                title='Upload new objects'
                onClose={this.props.onClose}>
                <form onSubmit={this._onSubmit}>
                    {this.state.loading ? 
                        <div className='ModalObjectAdd-progress'>
                            <ProgressBar
                                value={this.state.progress ? this.state.progress.loaded : 0}
                                total={this.state.progress ? this.state.progress.total : 1}
                            />
                        </div>
                    : null}

                    <div className='ModalObjectAdd-inputs'>
                        <InputFile
                            label='Files'
                            name='files'
                            error={this.state.formErrors.files}
                            onChange={this._onInputChange}
                        />
                        
                        {this.state.form.files && this.state.form.files.length === 1 ?
                            <InputText
                                label='Filename'
                                name='filename'
                                value={this.state.form.filename}
                                error={this.state.formErrors.filename}
                                onChange={this._onInputChange}
                            />
                        : null}
                    </div>
                    
                    <div className='ModalObjectAdd-actions'>
                        <Button
                            type='submit'
                            disabled={this.state.loading}>
                            Upload
                        </Button>
                    </div>
                </form>
            </Modal>
        );
    }
}
