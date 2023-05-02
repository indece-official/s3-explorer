import * as React from 'react';
import { Modal } from '../../Shared/Modal/Modal';
import { Link } from '../../Shared/Link/Link';

import './ModalUpdateAvailable.css';


export interface ModalUpdateAvailableProps
{
    onClose:    ( ) => any;
}


export class ModalUpdateAvailable extends React.Component<ModalUpdateAvailableProps>
{
    public render ( )
    {
        return (
            <Modal
                title='Update available'
                onClose={this.props.onClose}>
                <div className='ModalUpdateAvailable-title'>
                    S3 Explorer v{(window as any).VERSION || '???'}
                </div>

                <div className='ModalUpdateAvailable-text'>
                    TODO
                </div>
            </Modal>
        );
    }
}
