$module('streamchild', function () {
  var $static = this;
  $static.SimpleFunction = function (props, children) {
    var $temp0;
    var $temp1;
    var counter;
    var value;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            counter = $t.box(0, $g.____testlib.basictypes.Integer);
            $current = 1;
            continue;

          case 1:
            $temp1 = children;
            $current = 2;
            continue;

          case 2:
            $temp1.Next().then(function ($result0) {
              $temp0 = $result0;
              $result = $temp0;
              $current = 3;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 3:
            $result;
            value = $temp0.First;
            if ($t.unbox($temp0.Second)) {
              $current = 4;
              continue;
            } else {
              $current = 6;
              continue;
            }
            break;

          case 4:
            $g.____testlib.basictypes.Integer.$plus(counter, value).then(function ($result0) {
              counter = $result0;
              $result = counter;
              $current = 5;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 5:
            $result;
            $current = 2;
            continue;

          case 6:
            $g.____testlib.basictypes.Integer.$equals(counter, $t.box(9, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $result = $result0;
              $current = 7;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 7:
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
  $static.GetValues = function () {
    var $current = 0;
    var $continue = function ($yield, $yieldin, $reject, $done) {
      while (true) {
        switch ($current) {
          case 0:
            $yield($t.box(1, $g.____testlib.basictypes.Integer));
            $current = 1;
            return;

          case 1:
            $yield($t.box(2, $g.____testlib.basictypes.Integer));
            $current = 2;
            return;

          case 2:
            $yield($t.box(3, $g.____testlib.basictypes.Integer));
            $current = 3;
            return;

          default:
            $done();
            return;
        }
      }
    };
    return $generator.new($continue);
  };
  $static.TEST = function () {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.____testlib.basictypes.Mapping($g.____testlib.basictypes.String).Empty().then(function ($result1) {
              return $g.streamchild.GetValues().then(function ($result3) {
                return $g.____testlib.basictypes.MapStream($g.____testlib.basictypes.Integer, $g.____testlib.basictypes.Integer)($result3, function (value) {
                  var $current = 0;
                  var $continue = function ($resolve, $reject) {
                    while (true) {
                      switch ($current) {
                        case 0:
                          $g.____testlib.basictypes.Integer.$plus(value, $t.box(1, $g.____testlib.basictypes.Integer)).then(function ($result0) {
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
                }).then(function ($result2) {
                  return $g.streamchild.SimpleFunction($result1, $result2).then(function ($result0) {
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
