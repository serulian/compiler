$module('decorator', function () {
  var $static = this;
  $static.SimpleFunction = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $resolve($t.box(10, $g.____testlib.basictypes.Integer));
      return;
    };
    return $promise.new($continue);
  };
  $static.First = function (decorated, value) {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.Integer.$plus(decorated, value).then(function ($result0) {
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
  $static.Second = function (decorated, value) {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.Integer.$minus(decorated, value).then(function ($result0) {
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
  $static.Check = function (decorated, value) {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $promise.resolve($t.unbox(value)).then(function ($result0) {
              return ($promise.shortcircuit($result0, true) || $g.____testlib.basictypes.Integer.$equals(decorated, $t.box(15, $g.____testlib.basictypes.Integer))).then(function ($result1) {
                $result = $t.box($result0 && $t.unbox($result1), $g.____testlib.basictypes.Boolean);
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
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.decorator.SimpleFunction().then(function ($result3) {
              return $g.decorator.First($result3, $t.box(10, $g.____testlib.basictypes.Integer)).then(function ($result2) {
                return $g.decorator.Second($result2, $t.box(5, $g.____testlib.basictypes.Integer)).then(function ($result1) {
                  return $g.decorator.Check($result1, true).then(function ($result0) {
                    $result = $result0;
                    $current = 1;
                    $continue($resolve, $reject);
                    return;
                  });
                });
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
