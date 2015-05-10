/*global $ */
/*jslint browser: true */

(function () {
    'use strict';

    function keysHandler(evt) {

        var KEY_ENTER       = 13,
            KEY_LEFT_ARROW  = 37,
            KEY_RIGHT_ARROW = 39;

        if (evt.ctrlKey) {
            switch (evt.keyCode) {
                case KEY_ENTER:
                    $(':submit').trigger('click');
                    break;
                case KEY_LEFT_ARROW:
                    window.location = $('#pagination-prev').attr('href');
                    break;
                case KEY_RIGHT_ARROW:
                    window.location = $('#pagination-next').attr('href');
                    break;
            }
        }
    }

    function initKbdShortcuts() {
        $(window).keydown(keysHandler);
    }

    $(document).ready(function () {

        initKbdShortcuts();
        $('textarea').autosize();

    });

}());
