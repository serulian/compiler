$module('intconditional', function () {
  var $static = this;
  $static.TEST = function () {
    var $result;
    var first;
    var second;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            first = $t.fastbox(10, $g.____testlib.basictypes.Integer);
            second = $t.fastbox(2, $g.____testlib.basictypes.Integer);
            $g.____testlib.basictypes.Integer.$compare(second, first).then(function ($result0) {
              $result = $result0.$wrapped <= 0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            if ($result) {
              $current = 2;
              continue;
            } else {
              $current = 3;
              continue;
            }
            break;

          case 2:
            $resolve($t.fastbox(true, $g.____testlib.basictypes.Boolean));
            return;

          case 3:
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
