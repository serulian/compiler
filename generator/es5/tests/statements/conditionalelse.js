$module('conditionalelse', function () {
  var $static = this;
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            if (false) {
              $current = 1;
              continue;
            } else {
              $current = 2;
              continue;
            }
            break;

          case 1:
            $resolve($t.box(false, $g.____testlib.basictypes.Boolean));
            return;

          case 2:
            $resolve($t.box(true, $g.____testlib.basictypes.Boolean));
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
