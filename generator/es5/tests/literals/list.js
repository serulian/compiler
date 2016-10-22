$module('list', function () {
  var $static = this;
  $static.TEST = function () {
    var $result;
    var l;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.List($t.struct).forArray([$t.fastbox(1, $g.____testlib.basictypes.Integer), $t.fastbox(2, $g.____testlib.basictypes.Integer), $t.fastbox(3, $g.____testlib.basictypes.Integer), $t.fastbox(true, $g.____testlib.basictypes.Boolean)]).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            l = $result;
            l.Count().then(function ($result0) {
              $result = $t.fastbox($result0.$wrapped == 4, $g.____testlib.basictypes.Boolean);
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $resolve($result);
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
