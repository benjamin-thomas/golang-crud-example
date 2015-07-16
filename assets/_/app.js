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

    function initDataTables() {
        $('.datatable').each(function(i, table) {

            var tfoots = table.getElementsByTagName('tfoot');
            if (tfoots.length > 0) {
                var tfoot = tfoots[0];
                var ths = tfoot.getElementsByTagName('th');
                var tlength = ths.length;
                for (var ii = 0; ii < tlength; ii++ ) {
                    console.log("ii", ths[ii]);
                    var th = ths[ii];
                    var inputs = th.getElementsByTagName('input');
                    if (inputs.length > 0) {
                        inputs[0].placeholder = inputs[0].placeholder + '__search';
                    }
                }
            }

            $(table).dataTable({
                serverSide: true,
                ajax: table.dataset.ajax,
                deferRender: true,
                processing: true,
                paging: true,
                initComplete: function(_settings, _json) {
                    if (i !== 0) {
                        return;
                    }
                    var search = document.querySelector('input[type=search]');
                    search.focus();
                }
            });

            var dt = $(table).DataTable(); // DataTable(), not dataTable()!
            dt.columns().every(function() {
                var that = this;
                $('input', this.footer()).on('keyup change', function () {
                    that.search(this.value).draw();
                } );
            });

        });
    }

    $(document).ready(function () {

        initKbdShortcuts();
        autosize($('textarea'));

        var $paginationTotal = $('#pagination-count');

        $paginationTotal.on('click/delete', function(e) {
            this.textContent = parseInt(this.textContent) - 1;
        });

        $('a[data-method="DELETE"]').click(function(e) {
            e.preventDefault();

            var that = this;

            ajaxLink(this.href, 'DELETE')
                .done(function(data) {
                    var removeId = that.dataset.removeId,
                        redirect = that.dataset.redirect;

                    if (removeId) {
                        $('#' + removeId).remove();
                    } else if (redirect) {
                        window.location = redirect;
                    }
                    $paginationTotal.trigger('click/delete');
                })
                .fail(function(data) {
                    alert(data.status + ': ' + data.responseText);
                });

        });

        $('form[data-redirect]').submit(function(evt) {
            evt.preventDefault();

            var method = this.dataset.method || this.method;
            var redirect = this.dataset.redirect;

            ajaxLink(this.action, method, $(this).serialize())
                .done(function(data) {
                    window.location = redirect;
                })
                .fail(function(data) {
                    alert('Failed ' + method + ' request!');
                });
        });

    });

    $(document).ready(initDataTables);

}());
