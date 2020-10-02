// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import React from 'react';
import PropTypes from 'prop-types';

import {makeStyleFromTheme} from 'mattermost-redux/utils/theme_utils';

import {Svgs} from '../../constants';

export default class PostTypeWhiteboard extends React.PureComponent {
    static propTypes = {

        /*
         * The post to render the message for.
         */
        post: PropTypes.object.isRequired,

        /**
         * Set to render post body compactly.
         */
        compactDisplay: PropTypes.bool,

        /**
         * Flags if the post_message_view is for the RHS (Reply).
         */
        isRHS: PropTypes.bool,

        /**
         * Set to display times using 24 hours.
         */
        useMilitaryTime: PropTypes.bool,

        /*
         * Logged in user's theme.
         */
        theme: PropTypes.object.isRequired,

        /*
         * Creator's name.
         */
        creatorName: PropTypes.string.isRequired,

        /*
         * Current Channel Id.
         */
        currentChannelId: PropTypes.string.isRequired,

        /*
         * Whether the post was sent from a bot. Used for backwards compatibility.
         */
        fromBot: PropTypes.bool.isRequired,

        actions: PropTypes.shape({
            startWhiteboard: PropTypes.func.isRequired,
        }).isRequired,
    };

    static defaultProps = {
        mentionKeys: [],
        compactDisplay: false,
        isRHS: false,
    };

    constructor(props) {
        super(props);

        this.state = {
        };
    }

    render() {
        const style = getStyle(this.props.theme);
        const post = this.props.post;
        const props = post.props || {};

        let preText;

        preText = 'I have shared a whiteboard';
        if (this.props.fromBot) {
            preText = `${this.props.creatorName} has shared a whiteboard`;
        }
        const content = (
            <a
                className='btn btn-lg btn-primary'
                style={style.button}
                rel='noopener noreferrer'
                target='_blank'
                href={props.whiteboard_link}
            >
                <i
                    style={style.buttonIcon}
                    dangerouslySetInnerHTML={{__html: Svgs.WHITEBOARD}}
                />
                {'JOIN WHITEBOARD'}
            </a>
        );

        const subtitle = (
            <span>
                {'Whiteboard ID : '}
                <a
                    rel='noopener noreferrer'
                    target='_blank'
                    href={props.whiteboard_link}
                >
                    {props.whiteboard_id}
                </a>
            </span>
        );

        const title = 'Virtual Whiteboard';

        return (
            <div className='attachment attachment--pretext'>
                <div className='attachment__thumb-pretext'>
                    {preText}
                </div>
                <div className='attachment__content'>
                    <div className='clearfix attachment__container'>
                        <h5
                            className='mt-1'
                            style={style.title}
                        >
                            {title}
                        </h5>
                        {subtitle}
                        <div>
                            <div style={style.body}>
                                {content}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

const getStyle = makeStyleFromTheme((theme) => {
    return {
        body: {
            overflowX: 'auto',
            overflowY: 'hidden',
            paddingRight: '5px',
            width: '100%',
        },
        title: {
            fontWeight: '600',
        },
        button: {
            fontFamily: 'Open Sans',
            fontSize: '12px',
            fontWeight: 'bold',
            letterSpacing: '1px',
            lineHeight: '19px',
            marginTop: '12px',
            borderRadius: '4px',
            color: theme.buttonColor,
        },
        buttonIcon: {
            paddingRight: '8px',
            fill: theme.buttonColor,
        },
        summary: {
            fontFamily: 'Open Sans',
            fontSize: '14px',
            fontWeight: '600',
            lineHeight: '26px',
            margin: '0',
            padding: '14px 0 0 0',
        },
        summaryItem: {
            fontFamily: 'Open Sans',
            fontSize: '14px',
            lineHeight: '26px',
        },
    };
});
