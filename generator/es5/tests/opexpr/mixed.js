$module('mixed', function () {
  var $static = this;
  $static.TEST = function () {
    var $result;
    var finalIndex;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            finalIndex = $t.fastbox(-2, $g.____testlib.basictypes.Integer);
            $promise.resolve(finalIndex.$wrapped >= 10).then(function ($result0) {
              $result = $t.fastbox($result0 || (finalIndex.$wrapped < 0), $g.____testlib.basictypes.Boolean);
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
