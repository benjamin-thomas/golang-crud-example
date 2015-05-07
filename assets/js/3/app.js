(function() {

    function initSaveShortcut() {
        $(window).keydown(function(event) {
            var KEY_ENTER = 13;
            if(event.ctrlKey && (event.which === KEY_ENTER)) {
                $(':submit').trigger('click');
                // event.preventDefault();
            }
        });
    }

    $(document).ready(function() {

        initSaveShortcut();
        $('textarea').autosize();

    });
})();
