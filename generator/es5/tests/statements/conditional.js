$module('conditional', function () {
  var $static = this;
  $static.TEST = function () {
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            if (true) {
              $state.current = 1;
              continue;
            } else {
              $state.current = 2;
              continue;
            }
            break;

          case 1:
            $state.resolve($t.box(true, $g.____testlib.basictypes.Boolean));
            return;

          case 2:
            $state.resolve($t.box(false, $g.____testlib.basictypes.Boolean));
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
