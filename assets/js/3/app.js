/*global $ */
/*jslint browser: true */

(function () {
    'use strict';


    function ajaxLink(url, method, data) {
        return $.ajax({
            url: url,
            type: method,
            data: data,
        });
    }

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

        $('a[data-method="DELETE"]').click(function(e) {
            e.preventDefault();

            var that = this;

            ajaxLink(this.href, 'DELETE')
                .done(function(data) {
                    var remove = that.dataset.remove,
                        redirect = that.dataset.redirect;

                    if (remove) {
                        $('#' + remove).remove();
                    } else if (redirect) {
                        window.location = redirect;
                    }
                })
                .fail(function(data) {
                    alert(data.status + ': ' + data.responseText);
                });

        });

        $('form').submit(function(evt) {
            evt.preventDefault();

            var method = this.dataset.method || this.method;
            var redirect = this.dataset.redirect;

            ajaxLink(this.action, method, $(this).serialize())
                .done(function(data) {
                    if (redirect) {
                        window.location = redirect;
                    }
                })
                .fail(function(data) {
                    alert('Failed ' + method + ' request!');
                });
        });

    });

}());
