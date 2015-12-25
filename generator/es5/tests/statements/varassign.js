$module('varassign', function () {
  var $static = this;
  $static.DoSomething = function () {
    var i;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            i = 2;
            i = 3;
            1234;
            $state.current = -1;
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
