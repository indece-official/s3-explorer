import Moment from 'moment';
import * as React from 'react';
import { Modal } from '../../Shared/Modal/Modal';
import { ObjectV1 } from '../../Shared/Service/ObjectService';

import './ModalObjectDetails.css';


export interface ModalObjectDetailsProps
{
    profileID:  number;
    bucketName: string;
    object:     ObjectV1;
    onClose:    ( ) => any;
    onError:    ( err: Error ) => any;
}


export class ModalObjectDetails extends React.Component<ModalObjectDetailsProps>
{
    public render ( )
    {
        return (
            <Modal
                title='Object Details'
                onClose={this.props.onClose}>
                <div className='ModalObjectDetails-title'>
                    {this.props.object.key}
                </div>
                
                <div className='ModalObjectDetails-subtitle'>
                    {this.props.object.size} B
                </div>

                <div className='ModalObjectDetails-attributes'>
                    <div className='ModalObjectDetails-attribute'>
                        <div className='ModalObjectDetails-attribute-label'>
                            Bucket:
                        </div>

                        <div className='ModalObjectDetails-attribute-value'>
                            {this.props.bucketName}
                        </div>
                    </div>
                    
                    <div className='ModalObjectDetails-attribute'>
                        <div className='ModalObjectDetails-attribute-label'>
                            Owner:
                        </div>

                        <div className='ModalObjectDetails-attribute-value'>
                            {this.props.object.owner_name} ({this.props.object.owner_id})
                        </div>
                    </div>
                    
                    <div className='ModalObjectDetails-attribute'>
                        <div className='ModalObjectDetails-attribute-label'>
                            Last modified:
                        </div>

                        <div className='ModalObjectDetails-attribute-value'>
                            {Moment(this.props.object.last_modified).format('YYYY-MM-DD HH:mm:ss')}
                        </div>
                    </div>
                </div>
            </Modal>
        );
    }
}
