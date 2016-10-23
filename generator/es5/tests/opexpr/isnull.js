$module('isnull', function () {
  var $static = this;
  $static.TEST = function (a) {
    var $result;
    var b;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            b = $t.fastbox(1234, $g.____testlib.basictypes.Integer);
            $promise.resolve(a == null).then(function ($result0) {
              $result = $t.fastbox($result0 && (b != null), $g.____testlib.basictypes.Boolean);
              $current = 1;
              $continue($resolve, $reject);
              return;
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
