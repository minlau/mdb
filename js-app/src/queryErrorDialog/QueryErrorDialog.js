import React, {Component} from "react";
import {Button, Classes, Dialog, Icon, Intent} from "@blueprintjs/core";
import {Tooltip} from "@blueprintjs/core/lib/cjs";

class QueryErrorDialog extends Component {

    constructor(props) {
        super(props);

        this.state = {
            errors: [],
            isOpen: false
        };

        this.handleClose = this.handleClose.bind(this);
    }

    static getDerivedStateFromProps(nextProps, prevState) {
        if (nextProps.errors !== prevState.errors) {
            return {
                errors: nextProps.errors
            }
        } else return null;
    }

    handleClose() {
        this.setState({isOpen: false});
    }

    render() {
        this.state.errors.sort(function (a, b) {
            let o1 = a.groupId;
            let o2 = b.groupId;

            if (o1 < o2) return -1;
            if (o1 > o2) return 1;

            return 0;
        });
        const errorsContent = this.state.errors.map(error => {
            return (
                <div key={error.groupId}>
                    {error.groupId}. Position: {error.Position}, Message: {error.Message}

                    <Tooltip content={<div style={{whiteSpace: 'pre'}}>{JSON.stringify(error, null, 4)}</div>}>
                        <Icon style={{paddingLeft: '4px'}} icon="info-sign" intent={Intent.PRIMARY}/>
                    </Tooltip>
                </div>
            )
        });

        return (
            <Dialog
                onClose={this.handleClose}
                title="Errors"
                canEscapeKeyClose={true}
                canOutsideClickClose={true}
                enforceFocus={true}
                isOpen={this.state.isOpen}
                usePortal={true}
            >
                <div className={Classes.DIALOG_BODY}>
                    {errorsContent}
                </div>
                <div className={Classes.DIALOG_FOOTER}>
                    <div className={Classes.DIALOG_FOOTER_ACTIONS}>
                        <Button onClick={this.handleClose}>Close</Button>
                    </div>
                </div>
            </Dialog>
        );
    }
}

export default QueryErrorDialog;