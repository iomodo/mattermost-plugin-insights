import React from 'react';
import PropTypes from 'prop-types';
import {FormattedMessage} from 'react-intl';
import Backstage from './backstage';
import {CSSTransition} from 'react-transition-group';


const ANIMATION_DURATION = 100;

const Root = ({visible, close}) => {
    if (!visible) {
        return null;
    }

    return (
            <CSSTransition
            in={visible}
            classNames='FullScreenModal'
            mountOnEnter={true}
            unmountOnExit={true}
            timeout={ANIMATION_DURATION}
            appear={true}
        >
            <div className='FullScreenModal'>
                <Backstage
                    onBack={close}
                />
            </div>
        </CSSTransition>
    );
};

Root.propTypes = {
    visible: PropTypes.bool.isRequired,
    close: PropTypes.func.isRequired,
};

export default Root;
