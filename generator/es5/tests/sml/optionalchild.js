$module('optionalchild', function () {
  var $static = this;
  $static.SimpleFunction = function (props, child) {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $resolve($t.fastbox(child == null, $g.____testlib.basictypes.Boolean));
      return;
    };
    return $promise.new($continue);
  };
  $static.TEST = function () {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.Mapping($g.____testlib.basictypes.String).Empty().then(function ($result1) {
              return $g.optionalchild.SimpleFunction($result1).then(function ($result0) {
                $result = $result0;
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
