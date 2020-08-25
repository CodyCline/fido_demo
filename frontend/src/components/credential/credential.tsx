import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { ReactComponent as Key } from '../../assets/key-solid.svg';
import { ReactComponent as Trash } from '../../assets/trash-solid.svg';
import './credential.css';


export const Credential = ({ name, lastUsed, useCount, onDelete }: any) => {
    const [isShown, setShown] = React.useState<boolean>(false);
    const toggleModal = () => {
        setShown(!isShown);
    }
    const onConfirm = () => {
        toggleModal();
        onDelete();
    }
    return (
        <React.Fragment>
            <ConfirmModal
                show={isShown}
                onCancel={toggleModal}
                onConfirm={onConfirm}
            />
            <div className="credential">
                <div className="side-bar">
                    <Key style={{ height: "40px", color: "#CCC" }} />
                </div>
                <div className="credential-info">
                    <h3>{name || "Default Credential"}</h3>
                    <p>Last used: {lastUsed || "Never"}</p>
                    <p>Times used: {useCount || 0}</p>
                </div>
                <div onClick={toggleModal} className="delete-section">
                    <Trash className="delete-btn" />
                </div>
            </div>
        </React.Fragment>
    )
}

const ConfirmModal = ({ show, onCancel, onConfirm }: any) => {
    return (
        show && ReactDOM.createPortal(
            <div className="confirm-modal-bg">
                <div className="confirm-modal">
                    <p>Are you sure you want to delete this credential? You will no longer be able to use to login</p>
                    <button className="modal-btn cancel" onClick={onCancel}>NO</button>
                    <button className="modal-btn confirm" onClick={onConfirm}>YES</button>
                </div>
            </div>,
            document.body
        )
    );
}