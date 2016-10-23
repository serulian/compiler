$module('intconditional', function () {
  var $static = this;
  $static.TEST = function () {
    var first;
    var second;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            first = $t.fastbox(10, $g.____testlib.basictypes.Integer);
            second = $t.fastbox(2, $g.____testlib.basictypes.Integer);
            if (second.$wrapped <= first.$wrapped) {
              $current = 1;
              continue;
            } else {
              $current = 2;
              continue;
            }
            break;

          case 1:
            $resolve($t.fastbox(true, $g.____testlib.basictypes.Boolean));
            return;

          case 2:
            $resolve($t.fastbox(false, $g.____testlib.basictypes.Boolean));
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
