$module('string', function () {
  var $static = this;
  $static.DoSomething = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            'hello world';
            "hi world";
            'single quote with "quoted"';
            "double quote with 'quoted'";
            "escaped \" quote";
            'escaped \' quote';
            $state.current = -1;
            $state.returnValue = null;
            $callback($state);
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
});
