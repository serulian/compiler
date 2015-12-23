$module('boolean', function () {
  var $static = this;
  $static.SomeFunction = function (first, second) {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            first && second;
            first || second;
            !first;
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
