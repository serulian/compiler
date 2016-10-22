$module('basic', function () {
  var $static = this;
  $static.TEST = function () {
    var $result;
    var value;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            value = $t.fastbox(2, $g.____testlib.basictypes.Integer);
            $g.____testlib.basictypes.Integer.$equals(value, $t.fastbox(2, $g.____testlib.basictypes.Integer)).then(function ($result1) {
              return $promise.resolve($result1.$wrapped).then(function ($result0) {
                $result = $result0 ? $t.fastbox(true, $g.____testlib.basictypes.Boolean) : $t.fastbox(false, $g.____testlib.basictypes.Boolean);
                $current = 1;
                $continue($resolve, $reject);
                return;
              });
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
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
